package handlers

import (
	"encoding/json"
	"html/template"

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

	chartData, err := h.dashboardService.GetChartData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data grafik")
	}

	// Recent tracking activity
	recentTracking, err := h.dashboardService.GetRecentTracking(10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat aktivitas terbaru")
	}

	// Overdue alerts
	overdueCount, err := h.dashboardService.CountOverdueOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data keterlambatan")
	}

	overdueOrders, err := h.dashboardService.GetOverdueOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat daftar keterlambatan")
	}

	// Serialize chart data as JSON arrays for the template
	revenueJSON, _ := json.Marshal(chartData.Revenue)
	orderJSON, _ := json.Marshal(chartData.OrderCounts)
	labelsJSON, _ := json.Marshal(chartData.Labels)

	// Format timestamps for relative time display using template-friendly approach
	return render(c, "dashboard/index", fiber.Map{
		"title":              "Dashboard",
		"stats":              stats,
		"revenueData":        template.JS(revenueJSON),
		"orderData":          template.JS(orderJSON),
		"labels":             template.JS(labelsJSON),
		"recentTracking":     recentTracking,
		"overdueCount":       overdueCount,
		"overdueOrders":      overdueOrders,
		"timeAgo":            services.TimeAgo,
		"statusIcon":         services.StatusIcon,
		"statusLabel":        services.StatusLabel,
	}, "layouts/main")
}
