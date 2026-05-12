package main

import (
	"learn/config"
	"learn/internal/model"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env agar tahu koneksi database
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Koneksi ke Database
	config.ConnectDB()

	log.Println("Starting database migration...")

	// 3. Jalankan Auto Migration
	err := config.DB.AutoMigrate(&model.User{}, &model.Ticket{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed successfully!")

	// 4. Jalankan Seeder
	log.Println("Starting database seeding...")
	config.SeedData(config.DB)
	log.Println("Database seeding completed successfully!")
}
