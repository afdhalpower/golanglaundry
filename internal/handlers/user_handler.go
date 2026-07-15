package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	users, total, err := h.service.GetAll(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat data user")
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	roles := h.service.GetRoles()

	return c.Render("users/index", fiber.Map{
		"title":      "Manajemen User",
		"users":      users,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"roles":      roles,
	}, "layouts/main")
}

func (h *UserHandler) New(c fiber.Ctx) error {
	roles := h.service.GetRoles()
	return c.Render("users/form", fiber.Map{
		"title": "Tambah User",
		"user":  nil,
		"roles": roles,
	}, "layouts/main")
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	user, err := h.service.Create(
		c.FormValue("name"),
		c.FormValue("email"),
		c.FormValue("password"),
		c.FormValue("role"),
	)
	if err != nil {
		roles := h.service.GetRoles()
		return c.Render("users/form", fiber.Map{
			"title": "Tambah User",
			"error": err.Error(),
			"roles": roles,
		}, "layouts/main")
	}
	_ = user
	return c.Redirect().To("/users")
}

func (h *UserHandler) Edit(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	user, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User tidak ditemukan")
	}

	roles := h.service.GetRoles()
	return c.Render("users/form", fiber.Map{
		"title": "Edit User",
		"user":  user,
		"roles": roles,
	}, "layouts/main")
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	isActive := c.FormValue("is_active") == "on"

	if err := h.service.Update(uint(id), c.FormValue("name"), c.FormValue("email"), c.FormValue("role"), isActive); err != nil {
		user, _ := h.service.GetByID(uint(id))
		roles := h.service.GetRoles()
		return c.Render("users/form", fiber.Map{
			"title": "Edit User",
			"user":  user,
			"roles": roles,
			"error": err.Error(),
		}, "layouts/main")
	}

	return c.Redirect().To("/users")
}

func (h *UserHandler) ResetPassword(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	password := c.FormValue("password")

	if len(password) < 6 {
		return c.Status(fiber.StatusBadRequest).SendString("Password minimal 6 karakter")
	}

	if err := h.service.ResetPassword(uint(id), password); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect().To("/users")
}

func (h *UserHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	h.service.Delete(uint(id))
	return c.Redirect().To("/users")
}
