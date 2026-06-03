package public

import (
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type EventsHandler struct {
	events   *services.EventService
	settings *services.SiteSettingsService
	render   *render.Renderer
}

func NewEventsHandler(events *services.EventService, settings *services.SiteSettingsService, render *render.Renderer) *EventsHandler {
	return &EventsHandler{events: events, settings: settings, render: render}
}

func (h *EventsHandler) List(c *gin.Context) {
	items, _, _ := h.events.List(100, 0, false)
	h.render.Page(c, 200, "public/layout", "public/events_content", mergeSite(h.settings, gin.H{
		"Title":     "Афиша",
		"PageTitle": "Афиша",
		"Events":    items,
		"Active":    "events",
	}))
}
