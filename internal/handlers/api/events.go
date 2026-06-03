package api

import (
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
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

func toEventInput(req EventRequest) (services.EventInput, error) {
	date, err := parseDate(req.Date)
	if err != nil {
		return services.EventInput{}, err
	}
	return services.EventInput{
		Title:       req.Title,
		Date:        date,
		City:        req.City,
		Description: req.Description,
		TicketURL:   req.TicketURL,
	}, nil
}

func toAnySlice[T any](items []T) []any {
	out := make([]any, len(items))
	for i, v := range items {
		out[i] = v
	}
	return out
}
