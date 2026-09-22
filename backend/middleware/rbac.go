package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthorizeRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak: Hak akses tidak terdefinisi",
			})
		}

		allowed := false
		for _, r := range roles {
			if strings.EqualFold(userRole, r) {
				allowed = true
				break
			}
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Anda tidak memiliki akses ke resource ini",
			})
		}

		return c.Next()
	}
}