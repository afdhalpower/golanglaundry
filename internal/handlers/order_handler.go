package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/helpers"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type OrderHandler struct {
	orderService    *services.OrderService
	customerService *services.CustomerService
	serviceService  *services.ServiceService
}

func NewOrderHandler(
	orderService *services.OrderService,
	customerService *services.CustomerService,
	serviceService *services.ServiceService,
) *OrderHandler {
	return &OrderHandler{
		orderService:    orderService,
		customerService: customerService,
		serviceService:  serviceService,
	}
}

func (h *OrderHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")
	search := c.Query("search", "")

	orders, total, err := h.orderService.GetAll(page, limit, status, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data pesanan")
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	statuses := h.orderService.GetStatusList()
	counts, _ := h.orderService.GetStatusCounts()

	return render(c, "orders/index", fiber.Map{
		"title":      "Pesanan",
		"orders":     orders,
		"status":     status,
		"search":     search,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"statuses":   statuses,
		"counts":     counts,
	}, "layouts/main")
}

func (h *OrderHandler) New(c fiber.Ctx) error {
	customers, _, err := h.customerService.GetAll(1, 1000, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data")
	}

	services, err := h.serviceService.GetActive()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data layanan")
	}

	return render(c, "orders/form", fiber.Map{
		"title":     "Tambah Pesanan",
		"order":     nil,
		"customers": customers,
		"services":  services,
	}, "layouts/main")
}

func (h *OrderHandler) Create(c fiber.Ctx) error {
	customerID, _ := strconv.ParseUint(c.FormValue("customer_id"), 10, 32)
	serviceID, _ := strconv.ParseUint(c.FormValue("service_id"), 10, 32)
	weight, _ := strconv.ParseFloat(c.FormValue("weight_kg"), 64)
	discount, _ := strconv.ParseFloat(c.FormValue("discount"), 64)
	extraCost, _ := strconv.ParseFloat(c.FormValue("extra_cost"), 64)

	order, err := h.orderService.Create(services.CreateOrderRequest{
		CustomerID: uint(customerID),
		UserID:     helpers.LogAndGetUserID(c),
		ServiceID:  uint(serviceID),
		WeightKg:   weight,
		Discount:   discount,
		ExtraCost:  extraCost,
		EntryDate:  c.FormValue("entry_date"),
		Notes:      c.FormValue("notes"),
	})
	if err != nil {
		customers, _, _ := h.customerService.GetAll(1, 1000, "")
		services, _ := h.serviceService.GetActive()
		return render(c, "orders/form", fiber.Map{
			"title":     "Tambah Pesanan",
			"customers": customers,
			"services":  services,
			"error":     err.Error(),
		}, "layouts/main")
	}

	return c.Redirect().To("/orders/" + strconv.FormatUint(uint64(order.ID), 10))
}

func (h *OrderHandler) Show(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	order, err := h.orderService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pesanan tidak ditemukan")
	}

	validNext := h.orderService.GetValidNextStatuses(order.Status)

	return render(c, "orders/show", fiber.Map{
		"title":     "Detail Pesanan",
		"order":     order,
		"validNext": validNext,
	}, "layouts/main")
}

func (h *OrderHandler) UpdateStatus(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	newStatus := c.FormValue("status")
	note := c.FormValue("note")
	userID := helpers.LogAndGetUserID(c)

	order, err := h.orderService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pesanan tidak ditemukan")
	}

	if err := h.orderService.UpdateStatus(uint(id), newStatus, note, userID); err != nil {
		validNext := h.orderService.GetValidNextStatuses(order.Status)
		return render(c, "orders/show", fiber.Map{
			"title":     "Detail Pesanan",
			"order":     order,
			"validNext": validNext,
			"error":     err.Error(),
		}, "layouts/main")
	}

	return c.Redirect().To("/orders/" + c.Params("id"))
}

func (h *OrderHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	h.orderService.Delete(uint(id))
	return c.Redirect().To("/orders")
}
