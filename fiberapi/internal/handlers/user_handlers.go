package handlers

import (
  "fiberapi/internal/core/ports"
  "fiberapi/internal/utils"
  fiber "github.com/gofiber/fiber/v2"
  
  "fmt"
  "net/http"
  "log"
)

type UserHandlers struct {
  userService ports.UserService
}

type registerRequest struct {
  Username string `json:"username"`
  Email    string `json:"email"`
  Password string `json:"password"`
}

type loginRequest struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}


var _ ports.UserHandlers = (*UserHandlers)(nil)

func NewUserHandlers(userService ports.UserService) *UserHandlers {
  return &UserHandlers{
	userService: userService,
  }
}

func (h *UserHandlers) Login(c *fiber.Ctx) error {
  var req loginRequest

  if err := c.BodyParser(&req); err != nil {
      log.Printf("Error parsing request body: %v", err)
      return c.Status(http.StatusBadRequest).JSON(fiber.Map{
          "error": "cannot parse JSON",
      })
  }

  if req.Email == "" || req.Password == "" {
      log.Printf("Error: missing fields in request body")
      return c.Status(http.StatusBadRequest).JSON(fiber.Map{
          "error": "email and password are required",
      })
  }

  success, err := h.userService.Login(req.Email, req.Password)
  fmt.Println("Email: ", req.Email, " Password:", req.Password )
  if err != nil {
      log.Printf("Error logging in: %v", err)
      return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
          "error": "invalid email or password",
      })
  }

  if !success {
      return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
          "error": "invalid email or password",
      })
  }

  token, err := utils.GenerateJWT(req.Email, "user")
  if err != nil {
      log.Printf("Error generating JWT: %v", err)
      return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
          "error": "failed to generate token",
      })
  }

  return c.Status(http.StatusOK).JSON(fiber.Map{
      "code": "200",
      "auth": token,
  })
}


func (h *UserHandlers) Register(c *fiber.Ctx) error {
  var req registerRequest

  if err := c.BodyParser(&req); err != nil {
      return c.Status(http.StatusBadRequest).JSON(fiber.Map{
          "error": "cannot parse JSON",
      })
  }

  if req.Username == "" || req.Email == "" || req.Password == "" {
      return c.Status(http.StatusBadRequest).JSON(fiber.Map{
          "error": "username, email, and password are required",
      })
  }
  fmt.Println("Username: ", req.Username, " Email: ", req.Email, " Password: ", req.Password)

  err := h.userService.Register(req.Username, req.Email, req.Password, req.Password)
  if err != nil {
      if err.Error() == "user already exists" {
          return c.Status(http.StatusConflict).JSON(fiber.Map{
              "error": "user already exists",
          })
      }
      return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
          "error": "failed to register user",
      })
  }

  return c.Status(http.StatusCreated).JSON(fiber.Map{
      "message": "User added",
  })
}
