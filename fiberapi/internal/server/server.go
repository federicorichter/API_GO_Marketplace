package server

import (
  "fiberapi/internal/core/ports"
  "fiberapi/internal/middleware"
  "log"
  "os"

  fiber "github.com/gofiber/fiber/v2"
)

type Server struct {
  userHandlers ports.UserHandlers
  offerHandlers ports.OfferHandlers 
  adminHandlers ports.AdminHandlers
}

func NewServer(uHandlers ports.UserHandlers, oHandlers ports.OfferHandlers, aHandlers ports.AdminHandlers) *Server {
  return &Server{
      userHandlers:  uHandlers,
      offerHandlers: oHandlers,
      adminHandlers: aHandlers,
  }
}

func (s *Server) Initialize() {
  app := fiber.New()

  cppServer := app.Group("cpp_server")
  cppServer.Post("/update_supplies", s.offerHandlers.UpdateOffers)

  userRoutes := app.Group("/user")
  userRoutes.Post("/login", s.userHandlers.Login)
  userRoutes.Post("/register", s.userHandlers.Register)

  protectedRoutes := app.Group("/auth", middleware.JWTMiddleware())
  
  protectedRoutes.Get("/offers", s.offerHandlers.GetOffers)
  protectedRoutes.Post("/checkout", s.offerHandlers.Checkout)
  protectedRoutes.Get("/orders/:id", s.offerHandlers.GetStatus)

  app.Post("/loginadmin", s.adminHandlers.LoginAdmin)
  adminRoutes := app.Group("/admin", middleware.JWTMiddleware())
  adminRoutes.Use(middleware.JWTMiddleware(), middleware.AdminMiddleware())

  adminRoutes.Get("/dashboard", s.adminHandlers.GetDashboard)
  adminRoutes.Patch("/orders/:id", s.adminHandlers.PatchStatus)

  err := app.Listen(os.Getenv("PORT_API"))
    if err != nil {
	  log.Fatal(err)
  }
}
