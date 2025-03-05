package handlers

import (
  "fiberapi/internal/core/ports"
  "fiberapi/internal/utils"
  fiber "github.com/gofiber/fiber/v2"
  
  //"fmt"
  "net/http"
  "log"
  "strconv"
)

type AdminHandlers struct {
	adminService ports.AdminService
}

type updateStatusRequest struct {
    Status string `json:"status"`
}

var _ ports.AdminHandlers = (*AdminHandlers)(nil)

func NewAdminHandlers(adminService ports.AdminService) *AdminHandlers {
  return &AdminHandlers{
	adminService: adminService,
  }
}


func (h *AdminHandlers) LoginAdmin(c *fiber.Ctx) error {
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

    success, err := h.adminService.LoginAdmin(req.Email, req.Password)
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

    token, err := utils.GenerateJWT(req.Email, "admin")
    if err != nil {
        log.Printf("Error generating JWT: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to generate token",
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "code": "200",
        "auth_admin": token,
    })
}

func (h *AdminHandlers) GetDashboard (c *fiber.Ctx) error {
	offers, orders, sum, err := h.adminService.GetDashboard()

	if err != nil {
        log.Printf("Error fetching dashbord: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to fetch dashbord",
        })
    }
    return c.JSON(fiber.Map{
        "code": "200",
        "message": fiber.Map{
            "offers":offers,
			"orders":orders,
			"balance":sum,
        },
    })
}

func (h *AdminHandlers) PatchStatus (c *fiber.Ctx) error {

	orderIDParam := c.Params("id")
    orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
        log.Fatalf("Error with strconv: %v", err)
    }
    var req updateStatusRequest

    if err := c.BodyParser(&req); err != nil {
        log.Printf("Error parsing request body: %v", err)
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse JSON",
        })
    }

    if req.Status == "" {
        log.Printf("Error: missing status in request body")
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "status is required",
        })
    }

    err = h.adminService.PatchStatus(orderID, req.Status)
    if err != nil {
        log.Printf("Error updating order status: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to update order status",
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "code": "200",
        "message": "order status updated successfully",
    })

}
