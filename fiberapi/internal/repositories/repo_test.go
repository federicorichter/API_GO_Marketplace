package repositories

import (
	"context"
	"errors"
	"testing"
	//"net/http"

	//"fiberapi/internal/core/ports"

	//"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	//"github.com/jarcoal/httpmock"
)

func TestUserRepository_Register(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo,_ := NewUserRepository(mock)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT exists \\(SELECT 1 FROM users_table WHERE email=\\$1\\)").
			WithArgs("john@example.com").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO users_table \\(username, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\)").
			WithArgs("john_doe", "john@example.com", pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Register("john_doe", "john@example.com", "securepassword123")
		assert.NoError(t, err)
	})

	t.Run("user already exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT exists \\(SELECT 1 FROM users_table WHERE email=\\$1\\)").
			WithArgs("john@example.com").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.Register("john_doe", "john@example.com", "securepassword123")
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("database error on checking user existence", func(t *testing.T) {
		mock.ExpectQuery("SELECT exists \\(SELECT 1 FROM users_table WHERE email=\\$1\\)").
			WithArgs("john@example.com").
			WillReturnError(errors.New("database error"))

		err := repo.Register("john_doe", "john@example.com", "securepassword123")
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})

	t.Run("database error on user insertion", func(t *testing.T) {
		mock.ExpectQuery("SELECT exists \\(SELECT 1 FROM users_table WHERE email=\\$1\\)").
			WithArgs("john@example.com").
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO users_table \\(username, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\)").
			WithArgs("john_doe", "john@example.com", pgxmock.AnyArg()).
			WillReturnError(errors.New("insert error"))

		err := repo.Register("john_doe", "john@example.com", "securepassword123")
		assert.Error(t, err)
		assert.Equal(t, "insert error", err.Error())
	})
}

func TestUserLogin(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo,_ := NewUserRepository(mock)

	t.Run("login successfully", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepassword123"), bcrypt.DefaultCost)

		mock.ExpectQuery("SELECT password FROM users_table WHERE email=\\$1").
			WithArgs("john@example.com").
			WillReturnRows(pgxmock.NewRows([]string{"password"}).AddRow(string(hashedPassword)))

		resp, err := repo.Login("john@example.com", "securepassword123")
		assert.NoError(t, err)
		assert.True(t, resp)
	})

	t.Run("login with invalid password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepassword123"), bcrypt.DefaultCost)

		mock.ExpectQuery("SELECT password FROM users_table WHERE email=\\$1").
			WithArgs("john@example.com").
			WillReturnRows(pgxmock.NewRows([]string{"password"}).AddRow(string(hashedPassword)))

		resp, err := repo.Login("john@example.com", "wrongpassword")
		assert.Error(t, err)
		assert.False(t, resp)
		assert.Equal(t, "invalid password", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT password FROM users_table WHERE email=\\$1").
			WithArgs("john@example.com").
			WillReturnError(errors.New("no rows in result set"))

		resp, err := repo.Login("john@example.com", "securepassword123")
		assert.Error(t, err)
		assert.False(t, resp)
		assert.Equal(t, "no rows in result set", err.Error())
	})
}
/*
func TestCheckout(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo,_ := NewUserRepository(mock)

	// Mock the HTTP client
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Test case: Successful checkout
	t.Run("successful checkout", func(t *testing.T) {
		order := ports.Order{
			Items: []ports.Item{
				{ProductID: 1, Quantity: 2},
			},
		}

		mock.ExpectBegin()

		// Mock the response for fetching offers
		mock.ExpectQuery("SELECT id, quantity, price FROM offers WHERE id = \\$1").
			WithArgs(order.Items[0].ProductID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "quantity", "price"}).AddRow(1, 10, 100))

		// Mock the response for inserting into items_order
		mock.ExpectExec("INSERT INTO items_order \\(id_order, id_product, quantity\\) VALUES \\(\\$1, \\$2, \\$3\\)").
			WithArgs( order.Items[0].ProductID, order.Items[0].Quantity).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// Mock the response for updating offers
		mock.ExpectExec("UPDATE offers SET quantity = quantity - \\$1 WHERE id = \\$2").
			WithArgs(order.Items[0].Quantity, order.Items[0].ProductID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// Mock the response for inserting into status_orders
		mock.ExpectExec("INSERT INTO status_orders \\(id_order, status, total\\) VALUES \\(\\$1, \\$2, \\$3\\)").
			WithArgs( "processing", 200).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectCommit()

		// Mock the response for getting updated supplies
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("fruits").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("meat").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("vegetables").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("water").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("analgesics").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("antibiotics").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))
		mock.ExpectQuery("SELECT quantity FROM offers WHERE name = \\$1").
			WithArgs("bandages").WillReturnRows(pgxmock.NewRows([]string{"quantity"}).AddRow(5))

		// Mock the HTTP POST request to the C++ server
		httpmock.RegisterResponder("POST", "http://localhost:9045/update",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "")
				return resp, nil
			})

		totalSum, status, err := repo.Checkout(order)
		assert.NoError(t, err)
		assert.Equal(t, 200, totalSum)
		assert.Equal(t, "processing", status)
	})

	// Test case: Insufficient quantity
	t.Run("insufficient quantity", func(t *testing.T) {
		order := ports.Order{
			Items: []ports.Item{
				{ProductID: 1, Quantity: 2},
			},
		}

		mock.ExpectBegin()

		// Mock the response for fetching offers
		mock.ExpectQuery("SELECT id, quantity, price FROM offers WHERE id = \\$1").
			WithArgs(order.Items[0].ProductID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "quantity", "price"}).AddRow(1, 1, 100))

		mock.ExpectRollback()

		totalSum, status, err := repo.Checkout(order)
		assert.Error(t, err)
		assert.Equal(t, 0, totalSum)
		assert.Equal(t, "", status)
	})

	// Test case: Error during transaction
	t.Run("error during transaction", func(t *testing.T) {
		order := ports.Order{
			Items: []ports.Item{
				{ProductID: 1, Quantity: 2},
			},
		}

		mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))

		totalSum, status, err := repo.Checkout(order)
		assert.Error(t, err)
		assert.Equal(t, 0, totalSum)
		assert.Equal(t, "", status)
	})
}
*/
func TestGetOffers(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo,_ := NewUserRepository(mock)

	t.Run("successful fetch", func(t *testing.T) {
		mockRows := pgxmock.NewRows([]string{"id", "name", "quantity", "price", "category"}).
			AddRow(1, "Apples", 100, 1, "food").
			AddRow(2, "Bandages", 50, 2, "medicine")

		mock.ExpectQuery("SELECT id, name, quantity, price, category FROM offers").
			WillReturnRows(mockRows)

		offers, err := repo.GetOffers()
		assert.NoError(t, err)
		assert.Len(t, offers, 2)
		assert.Equal(t, "Apples", offers[0].Name)
		assert.Equal(t, 20, offers[0].Quantity) // 20% of 100
		assert.Equal(t, "Bandages", offers[1].Name)
		assert.Equal(t, 10, offers[1].Quantity) // 20% of 50
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, quantity, price, category FROM offers").
			WillReturnError(errors.New("query error"))

		offers, err := repo.GetOffers()
		assert.Error(t, err)
		assert.Nil(t, offers)
	})

	t.Run("scan error", func(t *testing.T) {
		mockRows := pgxmock.NewRows([]string{"id", "name", "quantity", "price", "category"}).
			AddRow(1, "Apples", "invalid", 1.5, "food")

		mock.ExpectQuery("SELECT id, name, quantity, price, category FROM offers").
			WillReturnRows(mockRows)

		offers, err := repo.GetOffers()
		assert.Error(t, err)
		assert.Nil(t, offers)
	})
}

func TestGetStatus(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo,_ := NewUserRepository(mock)

	t.Run("successful fetch", func(t *testing.T) {
		mock.ExpectQuery("SELECT status FROM status_orders WHERE id_order = \\$1").
			WithArgs(1).
			WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow("processing"))

		status, err := repo.GetStatus(1)
		assert.NoError(t, err)
		assert.Equal(t, "processing", status)
	})

	t.Run("no rows found", func(t *testing.T) {
		mock.ExpectQuery("SELECT status FROM status_orders WHERE id_order = \\$1").
			WithArgs(2).
			WillReturnError(errors.New("no rows in result set"))

		status, err := repo.GetStatus(2)
		assert.Error(t, err)
		assert.Equal(t, " ", status)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT status FROM status_orders WHERE id_order = \\$1").
			WithArgs(3).
			WillReturnError(errors.New("query error"))

		status, err := repo.GetStatus(3)
		assert.Error(t, err)
		assert.Equal(t, " ", status)
	})
}


