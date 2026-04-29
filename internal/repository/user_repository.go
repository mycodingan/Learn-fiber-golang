package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository digunakan untuk menginisialisasi repository user.
// Ini adalah implementasi dari Dependency Injection agar layer service bisa menggunakan repository ini.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Create menyimpan data user baru ke dalam database.
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByEmail mencari satu user berdasarkan alamat email.
// Digunakan saat proses Login atau pengecekan email unik saat Register.
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindByID mencari user berdasarkan ID (UUID).
func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}
