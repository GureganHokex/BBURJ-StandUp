package public

import (
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type PhotosHandler struct {
	photos   *services.PhotoService
	settings *services.SiteSettingsService
	render   *render.Renderer
}

func NewPhotosHandler(photos *services.PhotoService, settings *services.SiteSettingsService, render *render.Renderer) *PhotosHandler {
	return &PhotosHandler{photos: photos, settings: settings, render: render}
}

func (h *PhotosHandler) List(c *gin.Context) {
	items, _, _ := h.photos.List(100, 0)
	h.render.Page(c, 200, "public/layout", "public/photos_content", mergeSite(h.settings, gin.H{
		"Title":     "Фото",
		"PageTitle": "Фото",
		"Photos":    items,
		"Active":    "photos",
	}))
}
