package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type SettingHandler struct {
	service *services.SettingService
}

func NewSettingHandler(service *services.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

func (h *SettingHandler) Index(c fiber.Ctx) error {
	settings, _ := h.service.GetAll()
	return c.Render("settings/index", fiber.Map{
		"title":    "Pengaturan",
		"settings": settings,
	}, "layouts/main")
}

func (h *SettingHandler) Update(c fiber.Ctx) error {
	settings := map[string]string{
		"laundry_name":    c.FormValue("laundry_name"),
		"laundry_address": c.FormValue("laundry_address"),
		"laundry_phone":   c.FormValue("laundry_phone"),
		"laundry_whatsapp": c.FormValue("laundry_whatsapp"),
		"laundry_open":    c.FormValue("laundry_open"),
	}

	if err := h.service.Update(settings); err != nil {
		return c.Render("settings/index", fiber.Map{
			"title": "Pengaturan",
			"error": "Gagal menyimpan pengaturan",
		}, "layouts/main")
	}

	return c.Redirect().To("/settings")
}
