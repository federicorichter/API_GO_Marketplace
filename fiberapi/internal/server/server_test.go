package server

import (
	"bytes"
	"encoding/json"
	"fiberapi/internal/core/ports"
	"fiberapi/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Mock UserHandlers for testing
type MockUserHandlers struct {
	ports.UserHandlers
}

func (m *MockUserHandlers) Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success"})
}

func (m *MockUserHandlers) Login(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "token": "mocked_jwt_token"})
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserHandlers := &MockUserHandlers{}
	mockOfferHandlers := ports.NewMockOfferHandlers(ctrl)
	mockAdminHandlers := ports.NewMockAdminHandlers(ctrl)

	server := NewServer(mockUserHandlers, mockOfferHandlers, mockAdminHandlers)
	app := fiber.New()

	// Initialize routes with the app
	server.InitializeRoutes(app)

	// Create a request to send to the register endpoint
	req := httptest.NewRequest("POST", "/user/register", bytes.NewBufferString(`{ "username": "john_doe", "email": "john@example.com", "password": "securepassword123" }`))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the response body
	var responseBody map[string]interface{}
	err1 := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err1)
	assert.Equal(t, "success", responseBody["status"])
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserHandlers := &MockUserHandlers{}
	mockOfferHandlers := ports.NewMockOfferHandlers(ctrl)
	mockAdminHandlers := ports.NewMockAdminHandlers(ctrl)

	server := NewServer(mockUserHandlers, mockOfferHandlers, mockAdminHandlers)
	app := fiber.New()

	// Initialize routes with the app
	server.InitializeRoutes(app)

	// Create a request to send to the login endpoint
	req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString(`{ "email": "john@example.com", "password": "securepassword123" }`))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "mocked_jwt_token", responseBody["token"])
}

type MockOfferHandlers struct {
    ports.OfferHandlers
}

func (m *MockOfferHandlers) GetOffers(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"offers": []string{"offer1", "offer2"}})
}

func TestGetOffers(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserHandlers := &MockUserHandlers{}
    mockOfferHandlers := &MockOfferHandlers{}
    mockAdminHandlers := ports.NewMockAdminHandlers(ctrl)

    server := NewServer(mockUserHandlers, mockOfferHandlers, mockAdminHandlers)
    app := fiber.New()

    // Initialize routes with the app
    server.InitializeRoutes(app)

    // Create a request to send to the get offers endpoint
    req := httptest.NewRequest("GET", "/auth/offers", nil)
    req.Header.Set("Authorization", "Bearer mocktoken")

    // Execute the request
    resp, err := app.Test(req, -1)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Verify the response body
    var responseBody map[string]interface{}
    err2:= json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err2)
    assert.Equal(t, []interface{}{"offer1", "offer2"}, responseBody["offers"])
}

func TestCheckout(t *testing.T){
	ctrl := gomock.NewController(t)
    defer ctrl.Finish()

	mockUserHandlers := &MockUserHandlers{}
    mockOfferHandlers := &MockOfferHandlers{}
    mockAdminHandlers := ports.NewMockAdminHandlers(ctrl)

    server := NewServer(mockUserHandlers, mockOfferHandlers, mockAdminHandlers)
    app := fiber.New()

    // Initialize routes with the app
    server.InitializeRoutes(app)

	// Create a request to send to the get offers endpoint
    req := httptest.NewRequest("POST", "/auth/checkout", bytes.NewBufferString(`
	{ "id": 1, items: [{"quantity":10,"product_id":1},{"quantity":5, "product_id":4},{"quantity":3,"product_id":2}]}`))
    
	req.Header.Set("Authorization", "Bearer mocktoken")
}

// Mocking JWT middleware for testing protected routes
func MockJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mocking the JWT middleware, you can add specific claims here if needed
		c.Locals("user", "mockedUser")
		return c.Next()
	}
}

func (s *Server) InitializeRoutes(app *fiber.App) {
	// Mock JWT middleware for testing
	app.Use("/auth", MockJWTMiddleware())

	cppServer := app.Group("cpp_server")
	cppServer.Post("/update_supplies", s.offerHandlers.UpdateOffers)

	userRoutes := app.Group("/user")
	userRoutes.Post("/login", s.userHandlers.Login)
	userRoutes.Post("/register", s.userHandlers.Register)

	protectedRoutes := app.Group("/auth")
	protectedRoutes.Get("/offers", s.offerHandlers.GetOffers)
	protectedRoutes.Post("/checkout", s.offerHandlers.Checkout)
	protectedRoutes.Get("/orders/:id", s.offerHandlers.GetStatus)

	app.Post("/login_admin", s.adminHandlers.LoginAdmin)
	adminRoutes := app.Group("/admin")
	adminRoutes.Use(MockJWTMiddleware(), middleware.AdminMiddleware())

	adminRoutes.Get("/dashbord", s.adminHandlers.GetDashboard)
	adminRoutes.Patch("/orders/:id", s.adminHandlers.PatchStatus)
}
