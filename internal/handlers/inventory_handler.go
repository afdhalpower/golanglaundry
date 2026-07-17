package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/helpers"
	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	items, total, err := h.service.GetAll(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat inventaris")
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	lowStock, _ := h.service.GetLowStock()

	return render(c, "inventory/index", fiber.Map{
		"title":      "Inventaris",
		"items":      items,
		"search":     search,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"lowStock":   lowStock,
	}, "layouts/main")
}

func (h *InventoryHandler) New(c fiber.Ctx) error {
	return render(c, "inventory/form", fiber.Map{
		"title": "Tambah Barang",
		"item":  nil,
	}, "layouts/main")
}

func (h *InventoryHandler) Create(c fiber.Ctx) error {
	stock, _ := strconv.Atoi(c.FormValue("stock"))
	minStock, _ := strconv.Atoi(c.FormValue("min_stock"))

	item := &models.Inventory{
		Name:        c.FormValue("name"),
		Category:    c.FormValue("category"),
		Stock:       stock,
		MinStock:    minStock,
		Unit:        c.FormValue("unit"),
		Description: c.FormValue("description"),
	}

	userID := helpers.LogAndGetUserID(c)

	if err := h.service.Create(item, userID); err != nil {
		return render(c, "inventory/form", fiber.Map{
			"title": "Tambah Barang",
			"item":  item,
			"error": "Gagal menyimpan barang",
		}, "layouts/main")
	}

	return c.Redirect().To("/inventory")
}

func (h *InventoryHandler) Edit(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	item, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Barang tidak ditemukan")
	}

	return render(c, "inventory/form", fiber.Map{
		"title": "Edit Barang",
		"item":  item,
	}, "layouts/main")
}

func (h *InventoryHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	item, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Barang tidak ditemukan")
	}

	stock, _ := strconv.Atoi(c.FormValue("stock"))
	minStock, _ := strconv.Atoi(c.FormValue("min_stock"))

	item.Name = c.FormValue("name")
	item.Category = c.FormValue("category")
	item.Stock = stock
	item.MinStock = minStock
	item.Unit = c.FormValue("unit")
	item.Description = c.FormValue("description")

	userID := helpers.LogAndGetUserID(c)

	if err := h.service.Update(item, userID); err != nil {
		return render(c, "inventory/form", fiber.Map{
			"title": "Edit Barang",
			"item":  item,
			"error": "Gagal memperbarui barang",
		}, "layouts/main")
	}

	return c.Redirect().To("/inventory")
}

func (h *InventoryHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	h.service.Delete(uint(id))
	return c.Redirect().To("/inventory")
}

func (h *InventoryHandler) Movements(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	item, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Barang tidak ditemukan")
	}

	movements, err := h.service.GetMovements(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat riwayat stok")
	}

	return render(c, "inventory/movements", fiber.Map{
		"title":     "Riwayat Stok - " + item.Name,
		"item":      item,
		"movements": movements,
	}, "layouts/main")
}
