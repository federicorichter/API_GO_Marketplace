package handlers

import (
	"fiberapi/internal/core/ports"
    "github.com/gofiber/fiber/v2"
    "log"
	"strconv"
)

type SupplyUpdateRequest struct {
    Food     map[string]string `json:"food"`
    Medicine map[string]string `json:"medicine"`
}

type Food struct {
    Fruits     string `json:"fruits"`
    Meat       string `json:"meat"`
    Vegetables string `json:"vegetables"`
    Water      string `json:"water"`
}

type Medicine struct {
    Analgesics  string `json:"analgesics"`
    Antibiotics string `json:"antibiotics"`
    Bandages    string `json:"bandages"`
}

type UpdateRequest struct {
    Supplies SupplyUpdateRequest `json:"supplies"`
}

type OfferHandlers struct {
    userService ports.UserService
}

var _ ports.OfferHandlers = (*OfferHandlers)(nil) 

func NewOfferHandlers(userService ports.UserService) *OfferHandlers {
    return &OfferHandlers{
        userService: userService,
    }
}

func (h *OfferHandlers) GetOffers(c *fiber.Ctx) error {
    offers, err := h.userService.GetOffers()
    if err != nil {
        log.Printf("Error fetching offers: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to fetch offers",
        })
    }

    return c.JSON(offers)
}

func (h *OfferHandlers) Checkout(c *fiber.Ctx) error {
    var order ports.Order
    if err := c.BodyParser(&order); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
    }
	totalSum, status, err := h.userService.Checkout(order);
    if err != nil {
        log.Printf("Error in checkout: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to process order "})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "code": "200",
        "message": fiber.Map{
            "total":  totalSum,
            "status": status,
        },
    })
}

func(h *OfferHandlers) GetStatus(c *fiber.Ctx) error{
	orderIDParam := c.Params("id")
    orderID, err := strconv.Atoi(orderIDParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid order ID",
        })
    }

	status, err := h.userService.GetStatus(orderID)

	if err != nil {
        log.Printf("Error fetching order status: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to get order status",
        })
    }
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":"200",
		"message":fiber.Map{"status":status,},
	})
}

func (h *OfferHandlers) UpdateOffers(c *fiber.Ctx) error {
    var req UpdateRequest

    if err := c.BodyParser(&req); err != nil {
        log.Printf("Error parsing request body: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse JSON",
        })
    }

    err := h.userService.UpdateOffers(req.Supplies.Food, req.Supplies.Medicine)
    if err != nil {
        log.Printf("Error updating supplies: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to update supplies",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "code":    "200",
        "message": "supplies updated successfully",
    })
}

