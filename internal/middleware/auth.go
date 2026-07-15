package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var Store *session.Store

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		sess, err := Store.Get(c)
		if err != nil {
			return c.Redirect().To("/auth/login")
		}

		userID := sess.Get("user_id")
		if userID == nil {
			return c.Redirect().To("/auth/login")
		}

		c.Locals("user_id", userID)
		c.Locals("user_name", sess.Get("user_name"))
		c.Locals("user_role", sess.Get("user_role"))
		c.Locals("user_email", sess.Get("user_email"))

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
