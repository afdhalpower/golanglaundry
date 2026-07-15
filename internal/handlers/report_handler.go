package handlers

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func parseDateRange(c fiber.Ctx) (time.Time, time.Time) {
	now := time.Now()
	startStr := c.Query("start", "")
	endStr := c.Query("end", "")

	var start, end time.Time

	if startStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr != "" {
		end, _ = time.Parse("2006-01-02", endStr)
	}

	if start.IsZero() {
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	if end.IsZero() {
		end = now
	}

	end = end.Add(24*time.Hour - time.Second)

	return start, end
}

func (h *ReportHandler) Index(c fiber.Ctx) error {
	start, end := parseDateRange(c)

	revenue, daily, err := h.service.GetRevenueReport(start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat laporan pendapatan")
	}

	expense, err := h.service.GetExpenseReport(start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat laporan pengeluaran")
	}

	profit := revenue - expense

	orderCount, orderStats, err := h.service.GetOrderReport(start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat laporan order")
	}

	topCustomers, _ := h.service.GetTopCustomers()
	topServices, _ := h.service.GetTopServices()

	return c.Render("reports/index", fiber.Map{
		"title":        "Laporan",
		"startDate":    start.Format("2006-01-02"),
		"endDate":      end.Format("2006-01-02"),
		"totalRevenue": revenue,
		"totalExpense": expense,
		"profit":       profit,
		"orderCount":   orderCount,
		"orderStats":   orderStats,
		"revenueDaily": daily,
		"topCustomers": topCustomers,
		"topServices":  topServices,
	}, "layouts/main")
}

func (h *ReportHandler) ExportCSV(c fiber.Ctx) error {
	start, end := parseDateRange(c)
	reportType := c.Params("type", "revenue")

	c.Response().Header.Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=laporan_%s_%s.csv", reportType, time.Now().Format("20060102")))

	wr := csv.NewWriter(c.Response().BodyWriter())

	switch reportType {
	case "revenue":
		wr.Write([]string{"Tanggal", "Pendapatan", "Jumlah Transaksi"})
		wr.Write([]string{"Total", fmt.Sprintf("%.0f", 0.0), "0"})
	case "orders":
		wr.Write([]string{"Status", "Jumlah", "Total"})
		_, stats, _ := h.service.GetOrderReport(start, end)
		for _, s := range stats {
			wr.Write([]string{s.Status, fmt.Sprintf("%d", s.Count), fmt.Sprintf("%.0f", s.Total)})
		}
	case "profit":
		revenue, expense, profit, _ := h.service.GetProfitReport(start, end)
		wr.Write([]string{"Metrik", "Jumlah"})
		wr.Write([]string{"Pendapatan", fmt.Sprintf("%.0f", revenue)})
		wr.Write([]string{"Pengeluaran", fmt.Sprintf("%.0f", expense)})
		wr.Write([]string{"Laba", fmt.Sprintf("%.0f", profit)})
	}

	wr.Flush()
	return nil
}
