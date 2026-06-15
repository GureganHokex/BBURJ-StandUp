package api

import (
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/burj/comic/internal/tickets"
	"github.com/gin-gonic/gin"
)

type TicketCatalogHandler struct {
	catalog  *tickets.Catalog
	settings *services.SiteSettingsService
	events   *services.EventService
}

func NewTicketCatalogHandler(
	catalog *tickets.Catalog,
	settings *services.SiteSettingsService,
	events *services.EventService,
) *TicketCatalogHandler {
	return &TicketCatalogHandler{catalog: catalog, settings: settings, events: events}
}

func (h *TicketCatalogHandler) Providers(c *gin.Context) {
	cfg, err := h.ticketSettings()
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: h.catalog.Providers(cfg)})
}

func (h *TicketCatalogHandler) Events(c *gin.Context) {
	source := c.Param("source")
	cfg, err := h.ticketSettings()
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}

	items, err := h.catalog.ListEvents(c.Request.Context(), source, cfg)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: err.Error()})
		return
	}

	existing, err := h.events.ExternalIDs(source)
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	for i := range items {
		if _, ok := existing[items[i].ExternalID]; ok {
			items[i].AlreadyAdded = true
		}
	}

	c.JSON(http.StatusOK, ListResponse[any]{Data: toAnySlice(items), Meta: Meta{Total: int64(len(items))}})
}

func (h *TicketCatalogHandler) ticketSettings() (tickets.Settings, error) {
	s, err := h.settings.Get()
	if err != nil {
		return tickets.Settings{}, err
	}
	return tickets.SettingsFromModel(
		s.TimepadOrgID,
		s.TimepadAPIKey,
		s.TicketscloudOrgID,
		s.TicketscloudAPIKey,
		s.EventImportKeywords,
	), nil
}
