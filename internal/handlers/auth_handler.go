package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/afdhalpower/golanglaundry/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) LoginPage(c fiber.Ctx) error {
	m := session.FromContext(c)
	if m != nil && m.Get("user_id") != nil {
		return c.Redirect().To("/dashboard")
	}

	return render(c, "auth/login", fiber.Map{
		"title": "Login",
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return render(c, "auth/login", fiber.Map{
			"title":      "Login",
			"hideLayout": true,
			"error":      "Email dan password wajib diisi",
		}, "layouts/main")
	}

	user, err := h.authService.Login(services.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return render(c, "auth/login", fiber.Map{
			"title":      "Login",
			"hideLayout": true,
			"error":      err.Error(),
		}, "layouts/main")
	}

	m := session.FromContext(c)
	if m == nil {
		return render(c, "auth/login", fiber.Map{
			"title":      "Login",
			"hideLayout": true,
			"error":      "Gagal memulai session",
		}, "layouts/main")
	}

	m.Set("user_id", user.ID)
	m.Set("user_name", user.Name)
	m.Set("user_role", user.Role)
	m.Set("user_email", user.Email)

	return c.Redirect().To("/dashboard")
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	m := session.FromContext(c)
	if m != nil {
		m.Destroy()
	}
	return c.Redirect().To("/auth/login")
}
