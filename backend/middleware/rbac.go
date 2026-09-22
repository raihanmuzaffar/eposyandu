package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthorizeRoles memeriksa apakah role user sesuai dengan role yang diizinkan
func AuthorizeRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user")
		if userToken == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Akses ditolak: Token tidak ditemukan",
			})
		}

		token := userToken.(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userRole, ok := claims["role"].(string)

		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Akses ditolak: Role tidak valid",
			})
		}

		// Periksa apakah role user ada di dalam daftar allowedRoles
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next() // Role diizinkan, lanjutkan request
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Akses ditolak: Anda tidak memiliki hak akses untuk tindakan ini",
		})
	}
}