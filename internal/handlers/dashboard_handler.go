package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/services"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) Index(c fiber.Ctx) error {
	stats, err := h.dashboardService.GetStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat dashboard")
	}

	return render(c, "dashboard/index", fiber.Map{
		"title": "Dashboard",
		"stats": stats,
	}, "layouts/main")
}
