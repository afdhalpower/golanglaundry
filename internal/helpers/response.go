package helpers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSONSuccess(c fiber.Ctx, data interface{}) error {
	return c.JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

func JSONError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: message,
	})
}

func JSONValidationError(c fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(APIResponse{
		Success: false,
		Message: "Validasi gagal",
		Errors:  errors,
	})
}

func LogAndGetUserID(c fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		slog.Warn("user_id not found in context")
		return 0
	}
	return userID
}
