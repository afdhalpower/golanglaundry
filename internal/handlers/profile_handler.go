package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/helpers"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type ProfileHandler struct {
	authService *services.AuthService
}

func NewProfileHandler(authService *services.AuthService) *ProfileHandler {
	return &ProfileHandler{authService: authService}
}

func (h *ProfileHandler) Index(c fiber.Ctx) error {
	return c.Render("settings/profile", fiber.Map{
		"title": "Profil",
	}, "layouts/main")
}

func (h *ProfileHandler) Update(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	name := c.FormValue("name")
	email := c.FormValue("email")

	if name == "" || email == "" {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Nama dan email wajib diisi",
		}, "layouts/main")
	}

	if err := h.authService.UpdateProfile(userID, name, email); err != nil {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Gagal memperbarui profil",
		}, "layouts/main")
	}

	return c.Render("settings/profile", fiber.Map{
		"title":   "Profil",
		"success": "Profil berhasil diperbarui",
	}, "layouts/main")
}

func (h *ProfileHandler) ChangePassword(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")

	if newPassword == "" || len(newPassword) < 6 {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Password baru minimal 6 karakter",
		}, "layouts/main")
	}

	if err := h.authService.ChangePassword(userID, oldPassword, newPassword); err != nil {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": err.Error(),
		}, "layouts/main")
	}

	return c.Render("settings/profile", fiber.Map{
		"title":   "Profil",
		"success": "Password berhasil diubah",
	}, "layouts/main")
}
