package middleware

import (
	"learn/internal/dto"
	"learn/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware adalah middleware untuk mengecek apakah request memiliki token JWT yang valid.
// Jika valid, data user_id dan role akan disimpan di context (c.Locals) untuk digunakan oleh handler.
func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.WebResponse{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.WebResponse{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Invalid authorization format",
			})
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.WebResponse{
				Code:    fiber.StatusUnauthorized,
				Status:  "Unauthorized",
				Message: "Invalid or expired token",
			})
		}

		// Simpan data user dari token ke dalam context Fiber
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware digunakan untuk membatasi akses endpoint berdasarkan role user (misal: hanya 'admin').
func RoleMiddleware(requiredRole string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(dto.WebResponse{
				Code:    fiber.StatusForbidden,
				Status:  "Forbidden",
				Message: "You don't have permission to access this resource",
			})
		}
		return c.Next()
	}
}
