package handler

import (
	"learn/internal/dto"
	"learn/internal/service"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler menginisialisasi handler untuk endpoint autentikasi.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Register adalah endpoint POST /api/auth/register.
// Digunakan untuk mendaftarkan user baru ke sistem.
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
	}

	// Validasi dasar input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Name, Email and Password are required",
		})
	}

	res, err := h.authService.Register(req)
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
		Message: "User registered successfully",
		Data:    res,
	})
}

// Login adalah endpoint POST /api/auth/login.
// Digunakan untuk menukar email & password dengan JWT token.
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
	}

	res, err := h.authService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.WebResponse{
			Code:    fiber.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.WebResponse{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Login successful",
		Data:    res,
	})
}
