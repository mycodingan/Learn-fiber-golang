package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"learn/internal/model"
)

var DB *gorm.DB

// ConnectDB melakukan koneksi ke database MySQL menggunakan konfigurasi dari .env
func ConnectDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// 1. "Berteriak" ke MySQL untuk membuat database jika belum ada
	// Kita koneksi ke server MySQL tanpa memilih database terlebih dahulu
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	rawDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		log.Fatal("Gagal koneksi ke server MySQL:", err)
	}
	defer rawDB.Close()

	// Eksekusi perintah SQL untuk membuat database
	_, err = rawDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	if err != nil {
		log.Fatal("Gagal otomatis membuat database:", err)
	}

	log.Printf("Pastikan database '%s' siap digunakan", name)

	// 2. Koneksi utama menggunakan GORM setelah database dipastikan ada
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database MySQL connected successfully")

	// Auto migration
	err = DB.AutoMigrate(&model.User{}, &model.Ticket{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Jalankan Seeder
	SeedData(DB)
}
