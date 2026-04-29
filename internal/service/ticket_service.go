package service

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/repository"
)

type TicketService interface {
	Create(userID string, req dto.CreateTicketRequest) (*dto.TicketResponse, error)
	GetMyTickets(userID string, page, limit int) ([]dto.TicketResponse, int64, error)
	GetTicketByID(id string, userID string, role string) (*dto.TicketResponse, error)
	GetAllTickets(search string, page, limit int) ([]dto.TicketResponse, int64, error)
	UpdateStatus(id string, status model.TicketStatus) error
}

type ticketService struct {
	ticketRepo repository.TicketRepository
}

// NewTicketService menginisialisasi service untuk tiket.
func NewTicketService(ticketRepo repository.TicketRepository) TicketService {
	return &ticketService{ticketRepo}
}

// Create membuat tiket baru dengan status default 'open'.
func (s *ticketService) Create(userID string, req dto.CreateTicketRequest) (*dto.TicketResponse, error) {
	ticket := &model.Ticket{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
		Status:      model.StatusOpen,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	return &dto.TicketResponse{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		UserID:      ticket.UserID,
		CreatedAt:   ticket.CreatedAt.String(),
	}, nil
}

// GetMyTickets mengambil tiket-tiket yang dibuat oleh user yang sedang login.
func (s *ticketService) GetMyTickets(userID string, page, limit int) ([]dto.TicketResponse, int64, error) {
	offset := (page - 1) * limit
	tickets, total, err := s.ticketRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.TicketResponse
	for _, t := range tickets {
		res = append(res, dto.TicketResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			UserID:      t.UserID,
			CreatedAt:   t.CreatedAt.String(),
		})
	}
	return res, total, nil
}

// GetTicketByID mengambil detail satu tiket.
// Dilengkapi dengan pengecekan keamanan agar hanya pemilik tiket atau Admin yang bisa melihatnya.
func (s *ticketService) GetTicketByID(id string, userID string, role string) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Policy: Hanya pemilik tiket atau admin yang bisa melihat detail tiket
	if role != "admin" && ticket.UserID != userID {
		return nil, errors.New("forbidden")
	}

	return &dto.TicketResponse{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		UserID:      ticket.UserID,
		CreatedAt:   ticket.CreatedAt.String(),
	}, nil
}

// GetAllTickets mengambil semua tiket dari seluruh user (khusus Admin).
func (s *ticketService) GetAllTickets(search string, page, limit int) ([]dto.TicketResponse, int64, error) {
	offset := (page - 1) * limit
	tickets, total, err := s.ticketRepo.FindAll(search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.TicketResponse
	for _, t := range tickets {
		res = append(res, dto.TicketResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			UserID:      t.UserID,
			CreatedAt:   t.CreatedAt.String(),
		})
	}
	return res, total, nil
}

// UpdateStatus mengubah status progres tiket (khusus Admin).
func (s *ticketService) UpdateStatus(id string, status model.TicketStatus) error {
	_, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return errors.New("ticket not found")
	}

	return s.ticketRepo.UpdateStatus(id, status)
}
