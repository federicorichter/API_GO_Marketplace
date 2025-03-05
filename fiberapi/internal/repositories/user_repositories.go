package repositories

import (
	"context"
	"time"
	"errors"
	"fmt"
	"log"
    "bytes"
    "net/http"
    "encoding/json"

	"fiberapi/internal/core/ports"

    "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	PgClientTimeout = 5
)

type DBConn interface {
    QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
    Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    Begin(ctx context.Context) (pgx.Tx, error)
}


type UserRepository struct {
	conn DBConn
}

var _ ports.UserRepository = (*UserRepository)(nil)

func NewUserRepository(conn DBConn) (*UserRepository, error) {
	return &UserRepository{
		conn: conn,
	}, nil
}


func (r *UserRepository) Login(email, password string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var storedHashedPassword string
    err := r.conn.QueryRow(ctx, "SELECT password FROM users_table WHERE email=$1", email).Scan(&storedHashedPassword)
    if err != nil {
        if err == pgx.ErrNoRows {
            return false, errors.New("user not found")
        }
        return false, err
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
    if err != nil {
        return false, errors.New("invalid password")
    }

    return true, nil
}


func (r *UserRepository) Register(username, email, password string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if user already exists
    var exists bool
    err := r.conn.QueryRow(ctx, "SELECT exists (SELECT 1 FROM users_table WHERE email=$1)", email).Scan(&exists)
    if err != nil {
        log.Printf("Error checking if user exists: %v", err)
        return err
    }
    if exists {
        return errors.New("user already exists")
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return err
    }

    // Insert the user into the database
    _, err = r.conn.Exec(ctx, "INSERT INTO users_table (username, email, password) VALUES ($1, $2, $3)", username, email, string(hashedPassword))
    if err != nil {
        log.Printf("Error inserting user into database: %v", err)
        return err
    }

    return nil
}

func (r *UserRepository) GetOffers() ([]ports.Offer, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := r.conn.Query(ctx, "SELECT id, name, quantity, price, category FROM offers")
    if err != nil {
        log.Printf("Error fetching offers: %v", err)
        return nil, err
    }
    defer rows.Close()

    var offers []ports.Offer
    for rows.Next() {
        var offer ports.Offer
        err := rows.Scan(&offer.ID, &offer.Name, &offer.Quantity, &offer.Price, &offer.Category)
        if err != nil {
            log.Printf("Error scanning offer row: %v", err)
            return nil, err
        }
        offer.Quantity = int(offer.Quantity * 20/100) //only the 20% of supplies
        offers = append(offers, offer)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating offer rows: %v", err)
        return nil, err
    }

    return offers, nil
}

func (r *UserRepository) Checkout(order ports.Order) (int, string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    tx, err := r.conn.Begin(ctx)
    if err != nil {
        return 0, "", err
    }

    defer func() {
        if err != nil {
            _ = tx.Rollback(ctx)
        }
    }()
    
    var totalSum int

    for _, item := range order.Items {
        var id, quant, price int

        err := r.conn.QueryRow(ctx, "SELECT id, quantity, price FROM offers WHERE id = $1", item.ProductID).Scan(&id, &quant, &price)
        if err != nil {
            log.Printf("Error fetching offers: %v", err)
            return 0, "", err
        }

        if item.Quantity < 0 {
            return 0, "", fmt.Errorf("Quantity less than 0")
        }

        if int(quant*20/100) < item.Quantity {
            return 0, "", fmt.Errorf("insufficient quantity for product id: %d", item.ProductID)
        }

        totalSum += price * item.Quantity
        
    }

    for _, item := range order.Items {

        _, err = tx.Exec(ctx, "UPDATE offers SET quantity = quantity - $1 WHERE id = $2", item.Quantity, item.ProductID)
        if err != nil {
            log.Printf("Error updating offers: %v", err)
            return 0, "", err
        }
    }

    _, err = tx.Exec(ctx, "INSERT INTO status_orders (status, total) VALUES ($1, $2)", "processing", totalSum)
    if err != nil {
        log.Printf("Error inserting item into status_orders: %v", err)
        return 0, "", err
    }

    for _, item := range order.Items{

        var id int
        err := r.conn.QueryRow(ctx, "SELECT id_order FROM status_orders order by id_order limit 1").Scan(&id)
        if err != nil {
            log.Printf("Error fetching id: %v", err)
            return 0, "", err
        }

        _, err = tx.Exec(ctx, "INSERT INTO items_order (id_order, id_product, quantity) VALUES ($1, $2, $3)", id, item.ProductID, item.Quantity)
        if err != nil {
            log.Printf("Error inserting item into items_order: %v", err)
            return 0, "", err
        }
    }

    if err = tx.Commit(ctx); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return 0, "", err
    }

    updatedSupplies, err := r.getUpdatedSupplies(ctx)
    if err != nil {
        log.Printf("Error getting updated supplies: %v", err)
        return 0, "", err
    }

    err = sendUpdatedSuppliesToCppServer(updatedSupplies)
    if err != nil {
        log.Printf("Error sending updated supplies to C++ server: %v", err)
        return 0, "", err
    }

    return totalSum, "processing", nil
}


func (r *UserRepository) getUpdatedSupplies(ctx context.Context) (map[string]map[string]int, error) {
    supplies := make(map[string]map[string]int)
    supplies["food"] = make(map[string]int)
    supplies["medicine"] = make(map[string]int)

    foodItems := []string{"fruits", "meat", "vegetables", "water"}
    medicineItems := []string{"analgesics", "antibiotics", "bandages"}

    for _, item := range foodItems {
        var quantity int
        err := r.conn.QueryRow(ctx, "SELECT quantity FROM offers WHERE name = $1", item).Scan(&quantity)
        if err != nil {
            log.Printf("Error fetching offers from DB: %v", err)
            return nil, err
        }
        supplies["food"][item] = quantity
    }

    for _, item := range medicineItems {
        var quantity int
        
        err := r.conn.QueryRow(ctx, "SELECT quantity FROM offers WHERE name = $1", item).Scan(&quantity)
        
        if err != nil {
            return nil, err
        }
        supplies["medicine"][item] = quantity
    }

    return supplies, nil
}

func sendUpdatedSuppliesToCppServer(supplies map[string]map[string]int) error {
    jsonData, err := json.Marshal(supplies)
    if err != nil {
        log.Printf("Error parsing JSON: %v", err)
        return err
    }
    fmt.Println("GOT her" )
    req, err := http.NewRequest("POST", "http://localhost:9045/update", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error POSTING to C++ server: %v", err)
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error http client: %v", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to update C++ server: %s", resp.Status)
    }

    return nil
}

func (r *UserRepository) GetStatus(id int) (string, error){
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var status string

    err := r.conn.QueryRow(ctx, "SELECT status FROM status_orders WHERE id_order = $1", id).Scan(&status)
    if err != nil {
        log.Printf("Error fetching status from DB: %v", err)
        return " ", err
    }
    

    return status,nil
}

func (r *UserRepository) UpdateOffers(food, medicine map[string]string) error {
    fmt.Println("Food: ", food)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    tx, err := r.conn.Begin(ctx)
    if err != nil {
        return err
    }
    defer func() {
        err = tx.Rollback(ctx)
    }()
    if err != nil {
        return err
    }

    for item, quantity := range food {
        //_, quantityInt := strconv.Atoi(quantity)
        fmt.Println("Item: ", item, " Quant: ", quantity)
        _, err := tx.Exec(ctx, "UPDATE offers SET quantity = $1 WHERE category = 'food' AND name = $2", quantity, item)
        if err != nil {
            log.Printf("Error updating food item %s: %v", item, err)
            return err
        }
    }

    for item, quantity := range medicine {
        //_, quantityInt := strconv.Atoi(quantity)
        _, err := tx.Exec(ctx, "UPDATE offers SET quantity = $1 WHERE category = 'medicine' AND name = $2", quantity, item)
        if err != nil {
            log.Printf("Error updating medicine item %s: %v", item, err)
            return err
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return err
    }

    return nil
}
