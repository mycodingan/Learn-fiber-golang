package config

import (
	"learn/internal/model"
	"learn/utils"
	"log"

	"gorm.io/gorm"
)

// SeedData digunakan untuk mengisi data awal ke database.
func SeedData(db *gorm.DB) {
	// 1. Seed Admin
	var adminCount int64
	db.Model(&model.User{}).Where("role = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		hashedPassword, _ := utils.HashPassword("admin123")
		admin := model.User{
			Name:     "Super Admin",
			Email:    "admin@example.com",
			Password: hashedPassword,
			Role:     "admin",
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Println("Gagal seeding admin:", err)
		} else {
			log.Println("Admin berhasil di-seed: admin@example.com / admin123")
		}
	}

	// 2. Seed Regular User
	var userCount int64
	db.Model(&model.User{}).Where("role = ?", "user").Count(&userCount)

	if userCount == 0 {
		hashedPassword, _ := utils.HashPassword("user123")
		user := model.User{
			Name:     "Regular User",
			Email:    "user@example.com",
			Password: hashedPassword,
			Role:     "user",
		}
		if err := db.Create(&user).Error; err != nil {
			log.Println("Gagal seeding user:", err)
		} else {
			log.Println("User berhasil di-seed: user@example.com / user123")
		}
	}
}
