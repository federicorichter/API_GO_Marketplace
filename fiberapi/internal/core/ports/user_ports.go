package ports

import (

   "github.com/gofiber/fiber/v2"

)

type Offer struct {
  ID       int     `json:"id"`
  Name     string  `json:"name"`
  Quantity int     `json:"quantity"`
  Price    int `json:"price"`
  Category string  `json:"category"`
}

type Item struct {
  Quantity  int `json:"quantity"`
  ProductID int `json:"product_id"`
}

type Order struct {
  Items []Item `json:"items"`
}

type OrderStatus struct {
  ID int `json:"id"`
  Status string `json:"status"`
  Total int `json:"total"`
}

type UserRepository interface {
  Login( email, password string) (bool, error)
  Register(username string, email string, password string) error
  GetOffers() ([]Offer, error)
  Checkout(order Order) (int, string, error)
  GetStatus(id int) (string, error)
  UpdateOffers(food, medicine map[string]string) error
}

type UserService interface {
  Login(email, password string) (bool, error)
  Register(username string, email string, password string, passConfirm string) error
  GetOffers() ([]Offer, error)
  Checkout(order Order) (int, string, error)
  GetStatus(id int) (string, error)
  UpdateOffers(food, medicine map[string]string) error
}

type UserHandlers interface {
    Login(c *fiber.Ctx) error
    Register(c *fiber.Ctx) error
}

type OfferHandlers interface {
  GetOffers(c *fiber.Ctx) error
  Checkout(c *fiber.Ctx) error
  GetStatus(c *fiber.Ctx) error
  UpdateOffers(c *fiber.Ctx) error
}

