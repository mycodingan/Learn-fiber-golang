package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah teks password biasa menjadi hash terenkripsi menggunakan algoritma bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash membandingkan password teks biasa dengan hash yang tersimpan di database.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
