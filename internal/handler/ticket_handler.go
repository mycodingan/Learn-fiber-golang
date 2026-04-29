package handler

import (
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type TicketHandler struct {
	ticketService service.TicketService
}

// NewTicketHandler menginisialisasi handler untuk manajemen tiket.
func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService}
}

// Create menangani POST /api/tickets.
// Membuat tiket baru untuk user yang sedang login.
func (h *TicketHandler) Create(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req dto.CreateTicketRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
	}

	res, err := h.ticketService.Create(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.WebResponse{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Ticket created successfully",
		Data:    res,
	})
}

// GetMyTickets menangani GET /api/tickets.
// Menampilkan daftar tiket milik user sendiri.
func (h *TicketHandler) GetMyTickets(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.ticketService.GetMyTickets(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.WebResponse{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Tickets retrieved successfully",
		Data: fiber.Map{
			"tickets": res,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}

// GetTicketByID menangani GET /api/tickets/:id.
// Menampilkan detail satu tiket secara spesifik.
func (h *TicketHandler) GetTicketByID(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	res, err := h.ticketService.GetTicketByID(id, userID, role)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "ticket not found" {
			status = fiber.StatusNotFound
		} else if err.Error() == "forbidden" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(dto.WebResponse{
			Code:    status,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.WebResponse{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Ticket retrieved successfully",
		Data:    res,
	})
}

// GetAllTickets menangani GET /api/admin/tickets.
// Khusus Admin: Melihat semua tiket dari seluruh user.
func (h *TicketHandler) GetAllTickets(c fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	res, total, err := h.ticketService.GetAllTickets(search, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.WebResponse{
			Code:    fiber.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.WebResponse{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "All tickets retrieved successfully",
		Data: fiber.Map{
			"tickets": res,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}

// UpdateStatus menangani PUT /api/admin/tickets/:id.
// Khusus Admin: Memperbarui status pengerjaan tiket.
func (h *TicketHandler) UpdateStatus(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateTicketStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
	}

	// Validasi nilai status yang diperbolehkan
	if req.Status != model.StatusOpen && req.Status != model.StatusProcess && req.Status != model.StatusClosed {
		return c.Status(fiber.StatusBadRequest).JSON(dto.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid status value",
		})
	}

	err := h.ticketService.UpdateStatus(id, req.Status)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "ticket not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(dto.WebResponse{
			Code:    status,
			Status:  "Error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.WebResponse{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Ticket status updated successfully",
	})
}
