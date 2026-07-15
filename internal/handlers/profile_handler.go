package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/helpers"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type ProfileHandler struct {
	authService *services.AuthService
	userRepo    *repositories.UserRepository
}

func NewProfileHandler(authService *services.AuthService, userRepo *repositories.UserRepository) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (h *ProfileHandler) Index(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat profil")
	}

	return render(c, "settings/profile", fiber.Map{
		"title": "Profil",
		"user":  user,
	}, "layouts/main")
}

func (h *ProfileHandler) Update(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	name := c.FormValue("name")
	email := c.FormValue("email")

	if name == "" || email == "" {
		return render(c, "settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Nama dan email wajib diisi",
		}, "layouts/main")
	}

	if err := h.authService.UpdateProfile(userID, name, email); err != nil {
		return render(c, "settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Gagal memperbarui profil",
		}, "layouts/main")
	}

	// Reload user data
	user, _ := h.userRepo.FindByID(userID)
	return render(c, "settings/profile", fiber.Map{
		"title":   "Profil",
		"user":    user,
		"success": "Profil berhasil diperbarui",
	}, "layouts/main")
}

func (h *ProfileHandler) ChangePassword(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")

	if newPassword == "" || len(newPassword) < 6 {
		return render(c, "settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Password baru minimal 6 karakter",
		}, "layouts/main")
	}

	if err := h.authService.ChangePassword(userID, oldPassword, newPassword); err != nil {
		return render(c, "settings/profile", fiber.Map{
			"title": "Profil",
			"error": err.Error(),
		}, "layouts/main")
	}

	// Reload user data
	user, _ := h.userRepo.FindByID(userID)
	return render(c, "settings/profile", fiber.Map{
		"title":   "Profil",
		"user":    user,
		"success": "Password berhasil diubah",
	}, "layouts/main")
}
