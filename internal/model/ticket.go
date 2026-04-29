package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketStatus string

const (
	StatusOpen    TicketStatus = "open"
	StatusProcess TicketStatus = "process"
	StatusClosed  TicketStatus = "closed"
)

type Ticket struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TicketStatus   `gorm:"type:varchar(20);default:'open'" json:"status"`
	UserID      string         `gorm:"type:varchar(36);not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate adalah GORM Hook yang berjalan otomatis sebelum data Ticket dimasukkan ke database.
// Fungsi ini digunakan untuk men-generate UUID unik sebagai ID tiket.
func (t *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
