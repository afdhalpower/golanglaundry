package handlers

import "github.com/gofiber/fiber/v3"

// render wraps c.Render and injects user data (name, role) from c.Locals into template bindings.
// Every handler should use this instead of calling c.Render directly.
func render(c fiber.Ctx, name string, bind fiber.Map, layouts ...string) error {
	if bind == nil {
		bind = fiber.Map{}
	}
	if n := c.Locals("user_name"); n != nil {
		if _, ok := bind["name"]; !ok {
			bind["name"] = n
		}
	}
	if r := c.Locals("user_role"); r != nil {
		if _, ok := bind["role"]; !ok {
			bind["role"] = r
		}
	}
	return c.Render(name, bind, layouts...)
}
