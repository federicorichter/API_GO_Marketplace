package main

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	app := fiber.New()

	app.Post("/user/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success"})
	})

	req := httptest.NewRequest("POST", "/user/login", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRegister(t *testing.T) {
    app := fiber.New()

    app.Post("/user/register", func(c *fiber.Ctx) error {
        // Assume some registration logic here
        return c.JSON(fiber.Map{"status": "success"})
    })

    // Create a new HTTP request with the registration data
    req := httptest.NewRequest("POST", "/user/register", bytes.NewBufferString(`{ "email": "john@example.com", "password": "securepassword123" }`))
    req.Header.Set("Content-Type", "application/json")

    // Create a response recorder to capture the response
    resp, err := app.Test(req, -1)

    // Check for errors
    assert.NoError(t, err)

    // Verify the status code
    assert.Equal(t, http.StatusOK, resp.StatusCode)
/*
    // Optionally, verify the response body
    var responseBody map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&responseBody)
    assert.Equal(t, "success", responseBody["status"])*/
}

