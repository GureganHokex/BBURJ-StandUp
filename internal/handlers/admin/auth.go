package admin

import (
	"net/http"

	"github.com/burj/comic/internal/config"
	"github.com/burj/comic/internal/middleware"
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/security"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth   *services.AuthService
	csrf   *middleware.CSRF
	cfg    config.Config
	render *render.Renderer
}

func NewAuthHandler(auth *services.AuthService, csrf *middleware.CSRF, cfg config.Config, render *render.Renderer) *AuthHandler {
	return &AuthHandler{auth: auth, csrf: csrf, cfg: cfg, render: render}
}

func (h *AuthHandler) LoginPage(c *gin.Context) {
	sessionID := middleware.SessionID(c)
	if _, ok := h.auth.GetAdminID(sessionID); ok {
		c.Redirect(http.StatusFound, "/admin")
		return
	}
	h.render.HTML(c, 200, "admin/login.html", gin.H{
		"Title":    "Вход",
		"Next":     security.SafeRedirectPath(c.Query("next"), "/admin"),
		"CSRF":     h.csrf.TokenFromContext(c),
		"Error":    c.Query("error"),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	next := security.SafeRedirectPath(c.PostForm("next"), "/admin")

	sessionID, err := h.auth.Login(username, password)
	if err != nil {
		h.render.HTML(c, 401, "admin/login.html", gin.H{
			"Title": "Вход",
			"Next":  next,
			"CSRF":  h.csrf.TokenFromContext(c),
			"Error": "Неверный логин или пароль",
		})
		return
	}

	middleware.SetSessionCookie(c, h.cfg, sessionID)
	c.Redirect(http.StatusFound, next)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID := middleware.SessionID(c)
	h.auth.Logout(sessionID)
	middleware.ClearSessionCookie(c, h.cfg)
	c.Redirect(http.StatusFound, "/admin/login")
}
