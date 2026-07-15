package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/services"
	"github.com/afdhalpower/golanglaundry/internal/validation"
)

type ServiceHandler struct {
	service *services.ServiceService
}

func NewServiceHandler(service *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

func (h *ServiceHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	services, total, err := h.service.GetAll(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data layanan")
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return render(c, "services/index", fiber.Map{
		"title":      "Layanan",
		"services":   services,
		"search":     search,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	}, "layouts/main")
}

func (h *ServiceHandler) New(c fiber.Ctx) error {
	return render(c, "services/form", fiber.Map{
		"title":   "Tambah Layanan",
		"service": nil,
	}, "layouts/main")
}

func (h *ServiceHandler) Create(c fiber.Ctx) error {
	estimatedHours, _ := strconv.Atoi(c.FormValue("estimated_hours"))
	if estimatedHours <= 0 {
		estimatedHours = 24
	}

	price, _ := strconv.ParseFloat(c.FormValue("price_per_kg"), 64)

	service := &models.Service{
		Name:           c.FormValue("name"),
		PricePerKg:     price,
		EstimatedHours: estimatedHours,
		Description:    c.FormValue("description"),
		IsActive:       c.FormValue("is_active") == "on",
	}

	if errs := validation.ValidateStruct(service); errs != nil {
		return render(c, "services/form", fiber.Map{
			"title":   "Tambah Layanan",
			"service": service,
			"errors":  errs.ToMap(),
		}, "layouts/main")
	}

	if err := h.service.Create(service); err != nil {
		return render(c, "services/form", fiber.Map{
			"title":   "Tambah Layanan",
			"service": service,
			"error":   "Gagal menyimpan layanan",
		}, "layouts/main")
	}

	return c.Redirect().To("/services")
}

func (h *ServiceHandler) Edit(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	service, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Layanan tidak ditemukan")
	}

	return render(c, "services/form", fiber.Map{
		"title":   "Edit Layanan",
		"service": service,
	}, "layouts/main")
}

func (h *ServiceHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	service, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Layanan tidak ditemukan")
	}

	estimatedHours, _ := strconv.Atoi(c.FormValue("estimated_hours"))
	if estimatedHours <= 0 {
		estimatedHours = 24
	}

	price, _ := strconv.ParseFloat(c.FormValue("price_per_kg"), 64)

	service.Name = c.FormValue("name")
	service.PricePerKg = price
	service.EstimatedHours = estimatedHours
	service.Description = c.FormValue("description")
	service.IsActive = c.FormValue("is_active") == "on"

	if errs := validation.ValidateStruct(service); errs != nil {
		return render(c, "services/form", fiber.Map{
			"title":   "Edit Layanan",
			"service": service,
			"errors":  errs.ToMap(),
		}, "layouts/main")
	}

	if err := h.service.Update(service); err != nil {
		return render(c, "services/form", fiber.Map{
			"title":   "Edit Layanan",
			"service": service,
			"error":   "Gagal memperbarui layanan",
		}, "layouts/main")
	}

	return c.Redirect().To("/services")
}

func (h *ServiceHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal menghapus layanan")
	}
	return c.Redirect().To("/services")
}
