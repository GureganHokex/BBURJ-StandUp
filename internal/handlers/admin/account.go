package admin

import (
	"github.com/burj/comic/internal/middleware"
	"github.com/burj/comic/internal/render"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	csrf   *middleware.CSRF
	render *render.Renderer
}

func NewAccountHandler(csrf *middleware.CSRF, render *render.Renderer) *AccountHandler {
	return &AccountHandler{csrf: csrf, render: render}
}

func (h *AccountHandler) Page(c *gin.Context) {
	h.render.Page(c, 200, "admin/layout", "admin/account_content", gin.H{
		"Title":  "Смена пароля",
		"CSRF":   h.csrf.TokenFromContext(c),
		"Active": "account",
	})
}
