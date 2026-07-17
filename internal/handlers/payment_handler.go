package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/helpers"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")
	search := c.Query("search", "")

	payments, total, err := h.service.GetAll(page, limit, status, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat pembayaran")
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return render(c, "payments/index", fiber.Map{
		"title":      "Pembayaran",
		"payments":   payments,
		"status":     status,
		"search":     search,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	}, "layouts/main")
}

func (h *PaymentHandler) Create(c fiber.Ctx) error {
	orderID, _ := strconv.ParseUint(c.FormValue("order_id"), 10, 32)
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)

	_, err := h.service.CreateOrUpdate(
		uint(orderID),
		amount,
		c.FormValue("method"),
		c.FormValue("note"),
		helpers.LogAndGetUserID(c),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Redirect().To("/payments")
}

func (h *PaymentHandler) Show(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	payment, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pembayaran tidak ditemukan")
	}

	return render(c, "payments/show", fiber.Map{
		"title":   "Detail Pembayaran",
		"payment": payment,
	}, "layouts/main")
}
