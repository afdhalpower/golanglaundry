package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/afdhalpower/golanglaundry/internal/services"
)

type AuthHandler struct {
	authService  *services.AuthService
	sessionStore *session.Store
}

func NewAuthHandler(authService *services.AuthService, sessionStore *session.Store) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		sessionStore: sessionStore,
	}
}

func (h *AuthHandler) LoginPage(c fiber.Ctx) error {
	sess, err := h.sessionStore.Get(c)
	if err == nil {
		if sess.Get("user_id") != nil {
			return c.Redirect().To("/dashboard")
		}
	}

	return c.Render("auth/login", fiber.Map{
		"title":      "Login",
		"hideLayout": true,
	}, "layouts/main")
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Render("auth/login", fiber.Map{
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
		return c.Render("auth/login", fiber.Map{
			"title":      "Login",
			"hideLayout": true,
			"error":      err.Error(),
		}, "layouts/main")
	}

	sess, err := h.sessionStore.Get(c)
	if err != nil {
		return c.Render("auth/login", fiber.Map{
			"title":      "Login",
			"hideLayout": true,
			"error":      "Gagal memulai session",
		}, "layouts/main")
	}

	sess.Set("user_id", user.ID)
	sess.Set("user_name", user.Name)
	sess.Set("user_role", user.Role)
	sess.Set("user_email", user.Email)

	if err := sess.Save(); err != nil {
		return c.Render("auth/login", fiber.Map{
			"title":      "Login",
			"hideLayout": true,
			"error":      "Gagal menyimpan session",
		}, "layouts/main")
	}

	return c.Redirect().To("/dashboard")
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sess, err := h.sessionStore.Get(c)
	if err == nil {
		sess.Destroy()
	}
	return c.Redirect().To("/auth/login")
}
