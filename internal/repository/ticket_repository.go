package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ticket *model.Ticket) error
	FindAll(search string, limit, offset int) ([]model.Ticket, int64, error)
	FindByUserID(userID string, limit, offset int) ([]model.Ticket, int64, error)
	FindByID(id string) (*model.Ticket, error)
	UpdateStatus(id string, status model.TicketStatus) error
}

type ticketRepository struct {
	db *gorm.DB
}

// NewTicketRepository menginisialisasi repository untuk manajemen tiket.
func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

// Create membuat data tiket baru milik seorang user.
func (r *ticketRepository) Create(ticket *model.Ticket) error {
	return r.db.Create(ticket).Error
}

// FindAll mengambil semua tiket dengan fitur pencarian (search) dan pagination (limit, offset).
// Hanya bisa diakses oleh Admin.
func (r *ticketRepository) FindAll(search string, limit, offset int) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	query := r.db.Model(&model.Ticket{})
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Preload("User").Find(&tickets).Error

	return tickets, total, err
}

// FindByUserID mengambil daftar tiket yang dimiliki oleh user tertentu (berdasarkan userID).
func (r *ticketRepository) FindByUserID(userID string, limit, offset int) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	query := r.db.Model(&model.Ticket{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Find(&tickets).Error

	return tickets, total, err
}

// FindByID mencari detail satu tiket berdasarkan ID-nya.
func (r *ticketRepository) FindByID(id string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Preload("User").First(&ticket, "id = ?", id).Error
	return &ticket, err
}

// UpdateStatus mengubah status tiket (open, process, closed).
func (r *ticketRepository) UpdateStatus(id string, status model.TicketStatus) error {
	return r.db.Model(&model.Ticket{}).Where("id = ?", id).Update("status", status).Error
}
