package ports

import (

   "github.com/gofiber/fiber/v2"

)

type AdminRepository interface{
	LoginAdmin(email, password string) (bool, error)
	GetDashboard() ([]Offer, []OrderStatus, int, error)
	PatchStatus(id int, status string) error
  }
  
  type AdminService interface{
	LoginAdmin(email, password string) (bool, error)
	GetDashboard() ([]Offer, []OrderStatus, int, error)
	PatchStatus(id int, status string) error
  }
  
  type AdminHandlers interface{
	LoginAdmin(c *fiber.Ctx) error
	GetDashboard(c *fiber.Ctx) error
	PatchStatus(c *fiber.Ctx) error
  }
  