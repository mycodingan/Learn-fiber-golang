package routes

import (
	"learn/internal/handler"
	"learn/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, ticketHandler *handler.TicketHandler) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Ticket routes (User)
	tickets := api.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/", ticketHandler.GetMyTickets)
	tickets.Get("/:id", ticketHandler.GetTicketByID)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.Get("/tickets", ticketHandler.GetAllTickets)
	admin.Put("/tickets/:id", ticketHandler.UpdateStatus)
}
