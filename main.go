package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Inisialisasi Fiber app
	app := fiber.New()

	// Handler untuk route "/"
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Halo, Go Fiber v3 sudah berjalan!")
	})

	// Jalankan server di port 4567
	log.Fatal(app.Listen(":4567"))
}
