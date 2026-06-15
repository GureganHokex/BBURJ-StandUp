package public

import (
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	events   *services.EventService
	videos   *services.VideoService
	photos   *services.PhotoService
	merch    *services.MerchService
	settings *services.SiteSettingsService
	render   *render.Renderer
}

func NewHomeHandler(
	events *services.EventService,
	videos *services.VideoService,
	photos *services.PhotoService,
	merch *services.MerchService,
	settings *services.SiteSettingsService,
	render *render.Renderer,
) *HomeHandler {
	return &HomeHandler{events: events, videos: videos, photos: photos, merch: merch, settings: settings, render: render}
}

func (h *HomeHandler) Index(c *gin.Context) {
	upcoming, _, _ := h.events.List(4, 0, true)
	upcoming = services.NormalizeEventsForDisplay(upcoming)
	latest, _, _ := h.videos.List(1, 0)
	photos, _, _ := h.photos.List(12, 0)
	merchItems, _, _ := h.merch.List(4, 0)

	var featured *VideoView
	if len(latest) > 0 {
		v := latest[0]
		featured = &VideoView{Video: v, EmbedURL: services.EmbedURL(v.Platform, v.URL)}
	}

	h.render.Page(c, 200, "public/layout", "public/home_content", mergeSite(h.settings, gin.H{
		"Title":         "Главная",
		"PageTitle":     "Большой Буржинский",
		"Events":        upcoming,
		"FeaturedVideo": featured,
		"Photos":        photos,
		"Merch":         merchItems,
		"Active":        "home",
	}))
}
