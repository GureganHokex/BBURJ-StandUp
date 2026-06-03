package api

import (
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	service *services.PhotoService
}

func NewPhotoHandler(service *services.PhotoService) *PhotoHandler {
	return &PhotoHandler{service: service}
}

type PhotoRequest struct {
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

func (h *PhotoHandler) List(c *gin.Context) {
	limit, offset := parsePagination(c)
	items, total, err := h.service.List(limit, offset)
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[any]{Data: toAnySlice(items), Meta: Meta{Total: total}})
}

func (h *PhotoHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.service.Get(id)
	if writeNotFound(c, err) {
		return
	}
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: item})
}

func (h *PhotoHandler) Create(c *gin.Context) {
	var req PhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	item, errs, err := h.service.Create(services.PhotoInput{
		Title: req.Title, ImageURL: req.ImageURL, SortOrder: req.SortOrder,
	})
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	if errs != nil && errs.HasErrors() {
		writeValidationErrors(c, errs)
		return
	}
	c.JSON(http.StatusCreated, ItemResponse[any]{Data: item})
}

func (h *PhotoHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req PhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	item, errs, err := h.service.Update(id, services.PhotoInput{
		Title: req.Title, ImageURL: req.ImageURL, SortOrder: req.SortOrder,
	})
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
	c.JSON(http.StatusOK, ItemResponse[any]{Data: item})
}

func (h *PhotoHandler) Delete(c *gin.Context) {
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
