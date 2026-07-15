package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/services"
	"github.com/afdhalpower/golanglaundry/internal/validation"
)

type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	customers, total, err := h.service.GetAll(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data pelanggan")
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return render(c, "customers/index", fiber.Map{
		"title":      "Pelanggan",
		"customers":  customers,
		"search":     search,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	}, "layouts/main")
}

func (h *CustomerHandler) New(c fiber.Ctx) error {
	return render(c, "customers/form", fiber.Map{
		"title":    "Tambah Pelanggan",
		"customer": nil,
	}, "layouts/main")
}

func (h *CustomerHandler) Create(c fiber.Ctx) error {
	customer := &models.Customer{
		Name:     c.FormValue("name"),
		Phone:    c.FormValue("phone"),
		Whatsapp: c.FormValue("whatsapp"),
		Address:  c.FormValue("address"),
		Notes:    c.FormValue("notes"),
	}

	if errs := validation.ValidateStruct(customer); errs != nil {
		return render(c, "customers/form", fiber.Map{
			"title":    "Tambah Pelanggan",
			"customer": customer,
			"errors":   errs.ToMap(),
		}, "layouts/main")
	}

	if err := h.service.Create(customer); err != nil {
		return render(c, "customers/form", fiber.Map{
			"title":    "Tambah Pelanggan",
			"customer": customer,
			"error":    "Gagal menyimpan pelanggan",
		}, "layouts/main")
	}

	return c.Redirect().To("/customers")
}

func (h *CustomerHandler) Show(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	customer, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pelanggan tidak ditemukan")
	}

	return render(c, "customers/show", fiber.Map{
		"title":    "Detail Pelanggan",
		"customer": customer,
	}, "layouts/main")
}

func (h *CustomerHandler) Edit(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	customer, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pelanggan tidak ditemukan")
	}

	return render(c, "customers/form", fiber.Map{
		"title":    "Edit Pelanggan",
		"customer": customer,
	}, "layouts/main")
}

func (h *CustomerHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	customer, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pelanggan tidak ditemukan")
	}

	customer.Name = c.FormValue("name")
	customer.Phone = c.FormValue("phone")
	customer.Whatsapp = c.FormValue("whatsapp")
	customer.Address = c.FormValue("address")
	customer.Notes = c.FormValue("notes")

	if errs := validation.ValidateStruct(customer); errs != nil {
		return render(c, "customers/form", fiber.Map{
			"title":    "Edit Pelanggan",
			"customer": customer,
			"errors":   errs.ToMap(),
		}, "layouts/main")
	}

	if err := h.service.Update(customer); err != nil {
		return render(c, "customers/form", fiber.Map{
			"title":    "Edit Pelanggan",
			"customer": customer,
			"error":    "Gagal memperbarui pelanggan",
		}, "layouts/main")
	}

	return c.Redirect().To("/customers")
}

func (h *CustomerHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal menghapus pelanggan")
	}
	return c.Redirect().To("/customers")
}
