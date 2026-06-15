package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *services.EventService
	preview *services.URLPreviewService
}

func NewEventHandler(service *services.EventService, preview *services.URLPreviewService) *EventHandler {
	return &EventHandler{service: service, preview: preview}
}

func (h *EventHandler) List(c *gin.Context) {
	limit, offset := parsePagination(c)
	events, total, err := h.service.List(limit, offset, false)
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[any]{Data: toAnySlice(events), Meta: Meta{Total: total}})
}

func (h *EventHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	event, err := h.service.Get(id)
	if writeNotFound(c, err) {
		return
	}
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: event})
}

func (h *EventHandler) Create(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	input, err := toEventInput(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: map[string]string{"date": "invalid date format"}})
		return
	}
	event, errs, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateExternalEvent) {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "event already imported from this aggregator"})
			return
		}
		writeInternalError(c, appConfig(c), err)
		return
	}
	if errs != nil && errs.HasErrors() {
		writeValidationErrors(c, errs)
		return
	}
	c.JSON(http.StatusCreated, ItemResponse[any]{Data: event})
}

func (h *EventHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	input, err := toEventInput(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: map[string]string{"date": "invalid date format"}})
		return
	}
	event, errs, err := h.service.Update(id, input)
	if writeNotFound(c, err) {
		return
	}
	if err != nil {
		if errors.Is(err, services.ErrDuplicateExternalEvent) {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "event already imported from this aggregator"})
			return
		}
		writeInternalError(c, appConfig(c), err)
		return
	}
	if errs != nil && errs.HasErrors() {
		writeValidationErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: event})
}

func (h *EventHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *EventHandler) PreviewTicket(c *gin.Context) {
	rawURL := strings.TrimSpace(c.Query("url"))
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "url is required"})
		return
	}

	preview, err := h.preview.FetchPagePreview(c.Request.Context(), rawURL)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPreviewURLInvalid):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid url"})
		case errors.Is(err, services.ErrPreviewURLBlocked):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "url is not allowed"})
		case errors.Is(err, services.ErrPreviewNoImage):
			if previewHasData(preview) {
				c.JSON(http.StatusOK, ItemResponse[services.PagePreview]{Data: preview})
				return
			}
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "poster image not found on page"})
		default:
			writeInternalError(c, appConfig(c), err)
		}
		return
	}
	c.JSON(http.StatusOK, ItemResponse[services.PagePreview]{Data: preview})
}

func previewHasData(preview services.PagePreview) bool {
	return preview.Title != "" ||
		preview.Description != "" ||
		preview.City != "" ||
		preview.Date != "" ||
		preview.PosterImageURL != ""
}

func toEventInput(req EventRequest) (services.EventInput, error) {
	date, err := parseDate(req.Date)
	if err != nil {
		return services.EventInput{}, err
	}
	return services.EventInput{
		Title:          req.Title,
		Date:           date,
		City:           req.City,
		Description:    req.Description,
		TicketURL:      req.TicketURL,
		PosterImageURL: req.PosterImageURL,
		TicketSource:   req.TicketSource,
		ExternalID:     req.ExternalID,
	}, nil
}

func toAnySlice[T any](items []T) []any {
	out := make([]any, len(items))
	for i, v := range items {
		out[i] = v
	}
	return out
}
