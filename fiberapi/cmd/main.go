package main

import (
	"fiberapi/internal/core/services"
	"fiberapi/internal/handlers"
	"fiberapi/internal/repositories"
	"fiberapi/internal/server"
	"github.com/jackc/pgx/v4"

	"log"
	"context"
	"os"
)

func main() {
	
	dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }

	conn, err := pgx.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer conn.Close(context.Background())

	//repositories
    userRepository, err := repositories.NewUserRepository(conn)
    if err != nil {
        log.Fatalf("Error creating UserRepository: %v", err)
    }
	adminRepository, err := repositories.NewAdminRepository(dbURL)
	if err != nil {
        log.Fatalf("Error creating adminRepository: %v", err)
    }
	//services
	userService := services.NewUserService(userRepository)
	adminService := services.NewAdminService(adminRepository)
	//handlers
	offerHandlers := handlers.NewOfferHandlers(userService)
	userHandlers := handlers.NewUserHandlers(userService)
	adminHandlers := handlers.NewAdminHandlers(adminService)
	//server
	httpServer := server.NewServer(
		userHandlers,
		offerHandlers,
		adminHandlers,
	)
	httpServer.Initialize()
}
