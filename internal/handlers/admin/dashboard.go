package admin

import (
	"github.com/burj/comic/internal/admin"
	"github.com/burj/comic/internal/middleware"
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	events   *services.EventService
	videos   *services.VideoService
	merch    *services.MerchService
	photos   *services.PhotoService
	settings *services.SiteSettingsService
	csrf     *middleware.CSRF
	render   *render.Renderer
}

func NewDashboardHandler(
	events *services.EventService,
	videos *services.VideoService,
	merch *services.MerchService,
	photos *services.PhotoService,
	settings *services.SiteSettingsService,
	csrf *middleware.CSRF,
	render *render.Renderer,
) *DashboardHandler {
	return &DashboardHandler{
		events: events, videos: videos, merch: merch, photos: photos,
		settings: settings, csrf: csrf, render: render,
	}
}

type modelCard struct {
	admin.Model
	Count int64
}

func (h *DashboardHandler) Index(c *gin.Context) {
	cards := make([]modelCard, 0)
	for _, m := range admin.Models() {
		var count int64
		switch m.Slug {
		case "settings":
			count = 1
		case "events":
			count, _ = h.events.Count()
		case "videos":
			count, _ = h.videos.Count()
		case "photos":
			count, _ = h.photos.Count()
		case "merch":
			count, _ = h.merch.Count()
		}
		cards = append(cards, modelCard{Model: m, Count: count})
	}
	h.render.Page(c, 200, "admin/layout", "admin/index_content", gin.H{
		"Title":  "Админ-панель",
		"Cards":  cards,
		"CSRF":   h.csrf.TokenFromContext(c),
		"Active": "",
	})
}
