package api

import (
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	service *services.VideoService
}

func NewVideoHandler(service *services.VideoService) *VideoHandler {
	return &VideoHandler{service: service}
}

func (h *VideoHandler) List(c *gin.Context) {
	limit, offset := parsePagination(c)
	videos, total, err := h.service.List(limit, offset)
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[any]{Data: toAnySlice(videos), Meta: Meta{Total: total}})
}

func (h *VideoHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	video, err := h.service.Get(id)
	if writeNotFound(c, err) {
		return
	}
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: video})
}

func (h *VideoHandler) Create(c *gin.Context) {
	var req VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	video, errs, err := h.service.Create(services.VideoInput{Title: req.Title, URL: req.URL})
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	if errs != nil && errs.HasErrors() {
		writeValidationErrors(c, errs)
		return
	}
	c.JSON(http.StatusCreated, ItemResponse[any]{Data: video})
}

func (h *VideoHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	video, errs, err := h.service.Update(id, services.VideoInput{Title: req.Title, URL: req.URL})
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
	c.JSON(http.StatusOK, ItemResponse[any]{Data: video})
}

func (h *VideoHandler) Delete(c *gin.Context) {
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
