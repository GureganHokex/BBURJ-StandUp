package admin

import (
	"github.com/burj/comic/internal/middleware"
	"github.com/burj/comic/internal/render"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	csrf   *middleware.CSRF
	render *render.Renderer
}

func NewSettingsHandler(csrf *middleware.CSRF, render *render.Renderer) *SettingsHandler {
	return &SettingsHandler{csrf: csrf, render: render}
}

func (h *SettingsHandler) Page(c *gin.Context) {
	h.render.Page(c, 200, "admin/layout", "admin/settings_content", gin.H{
		"Title":  "Настройки сайта",
		"CSRF":   h.csrf.TokenFromContext(c),
		"Active": "settings",
	})
}
