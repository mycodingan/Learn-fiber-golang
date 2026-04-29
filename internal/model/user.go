package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"type:varchar(10);default:'user'" json:"role"` // 'admin' or 'user'
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Tickets   []Ticket       `gorm:"foreignKey:UserID" json:"tickets,omitempty"`
}

// BeforeCreate adalah GORM Hook yang berjalan otomatis sebelum data User dimasukkan ke database.
// Fungsi ini digunakan untuk men-generate UUID unik sebagai ID user.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
