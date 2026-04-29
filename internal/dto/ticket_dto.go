package dto

import "learn/internal/model"

type CreateTicketRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTicketStatusRequest struct {
	Status model.TicketStatus `json:"status" validate:"required,oneof=open process closed"`
}

type TicketResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      model.TicketStatus `json:"status"`
	UserID      string             `json:"user_id"`
	CreatedAt   string             `json:"created_at"`
}

type WebResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
