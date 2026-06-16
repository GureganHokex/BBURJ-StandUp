package public

import (
	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type MerchHandler struct {
	merch    *services.MerchService
	settings *services.SiteSettingsService
	render   *render.Renderer
}

func NewMerchHandler(merch *services.MerchService, settings *services.SiteSettingsService, render *render.Renderer) *MerchHandler {
	return &MerchHandler{merch: merch, settings: settings, render: render}
}

func (h *MerchHandler) List(c *gin.Context) {
	if !denyHiddenSection(c, h.settings, func(s *models.SiteSettings) bool { return s.ShowMerch }) {
		return
	}
	items, _, _ := h.merch.List(100, 0)
	h.render.Page(c, 200, "public/layout", "public/merch_content", mergeSite(h.settings, gin.H{
		"Title":     "Мерч",
		"PageTitle": "Мерч",
		"Merch":     items,
		"Active":    "merch",
	}))
}
