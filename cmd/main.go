package main

import (
	"learn/config"
	"learn/internal/handler"
	"learn/internal/repository"
	"learn/internal/service"
	"learn/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Connect to Database
	config.ConnectDB()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"code":    code,
				"status":  "Error",
				"message": err.Error(),
			})
		},
	})

	// Global Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize Repositories
	userRepo := repository.NewUserRepository(config.DB)
	ticketRepo := repository.NewTicketRepository(config.DB)

	// Initialize Services
	authService := service.NewAuthService(userRepo)
	ticketService := service.NewTicketService(ticketRepo)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	// Setup Routes
	routes.SetupRoutes(app, authHandler, ticketHandler)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
