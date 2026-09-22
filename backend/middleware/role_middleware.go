package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRoles membatasi akses endpoint berdasarkan role pengguna
func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role yang disimpan dari JWT Middleware
		userRole, ok := c.Locals("user_role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Akses tidak sah: Informasi role tidak ditemukan",
			})
		}

		// Periksa apakah role user termasuk dalam daftar role yang diperbolehkan
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Akses ditolak: Anda tidak memiliki hak akses untuk fitur ini",
		})
	}
}