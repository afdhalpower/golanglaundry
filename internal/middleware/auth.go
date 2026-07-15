package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		m := session.FromContext(c)
		if m == nil {
			return c.Redirect().To("/auth/login")
		}

		userID := m.Get("user_id")
		if userID == nil {
			return c.Redirect().To("/auth/login")
		}

		c.Locals("user_id", userID)
		c.Locals("user_name", m.Get("user_name"))
		c.Locals("user_role", m.Get("user_role"))
		c.Locals("user_email", m.Get("user_email"))

		return c.Next()
	}
}

func RolePermission(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("user_role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("Akses ditolak")
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).SendString("Akses ditolak")
	}
}
