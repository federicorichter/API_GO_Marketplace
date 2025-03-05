package repositories

import (
	"context"
	"time"
	"errors"
	"fmt"
	"log"

	"fiberapi/internal/core/ports"

	"github.com/jackc/pgx/v4"
	//"golang.org/x/crypto/bcrypt"
)

type AdminRepository struct {
	conn *pgx.Conn
}

var _ ports.AdminRepository = (*AdminRepository)(nil)

func NewAdminRepository(connString string) (*AdminRepository, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), PgClientTimeout*time.Second)
	defer cancelFunc()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &AdminRepository{
		conn: conn,
	}, nil
}

func (r *AdminRepository) LoginAdmin(email, password string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var storedHashedPassword string
    err := r.conn.QueryRow(ctx, "SELECT password FROM admins WHERE email=$1", email).Scan(&storedHashedPassword)
    if err != nil {
        if err == pgx.ErrNoRows {
            return false, errors.New("user not found")
        }
        return false, err
    }
    fmt.Println("EmailAdmin: ", email, " PasswordAdmin:", storedHashedPassword )
    //err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
    if err != nil {
        return false, errors.New("invalid password")
    }

    return true, nil
}

func (r *AdminRepository) GetDashboard() ([]ports.Offer, []ports.OrderStatus, int, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := r.conn.Query(ctx, "SELECT id, name, quantity, price, category FROM offers")
    if err != nil {
        log.Printf("Error fetching offers: %v", err)
        return nil, nil,  0, err
    }
    defer rows.Close()

    var offers []ports.Offer
    for rows.Next() {
        var offer ports.Offer
        err := rows.Scan(&offer.ID, &offer.Name, &offer.Quantity, &offer.Price, &offer.Category)
        if err != nil {
            log.Printf("Error scanning offer row: %v", err)
            return nil, nil,  0, err
        }
        offers = append(offers, offer)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating offer rows: %v", err)
        return nil, nil,  0, err
    }

	rows, err = r.conn.Query(ctx, "SELECT id_order, status, total FROM status_orders")
    if err != nil {
        log.Printf("Error fetching status_orders: %v", err)
        return nil, nil,  0, err
    }
    defer rows.Close()

	var orders []ports.OrderStatus
	var sum = 0
	for rows.Next() {
		var order ports.OrderStatus
		err = rows.Scan(&order.ID, &order.Status, &order.Total)
		sum += order.Total
		if err != nil {
            log.Printf("Error scanning offer row: %v", err)
            return nil, nil,  0, err
        }
        orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
        log.Printf("Error iterating status rows: %v", err)
        return nil, nil,  0, err
    }

	return offers, orders, sum, nil
}

func (r *AdminRepository) PatchStatus(id int, status string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := "UPDATE status_orders SET status = $1 WHERE id_order = $2"
    _, err := r.conn.Exec(ctx, query, status, id)
    if err != nil {
        log.Printf("Error updating order status in database: %v", err)
        return err
    }

    return nil

}
