package api

import (
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type MerchHandler struct {
	service *services.MerchService
}

func NewMerchHandler(service *services.MerchService) *MerchHandler {
	return &MerchHandler{service: service}
}

func (h *MerchHandler) List(c *gin.Context) {
	limit, offset := parsePagination(c)
	items, total, err := h.service.List(limit, offset)
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[any]{Data: toAnySlice(items), Meta: Meta{Total: total}})
}

func (h *MerchHandler) Get(c *gin.Context) {
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

func (h *MerchHandler) Create(c *gin.Context) {
	var req MerchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	item, errs, err := h.service.Create(services.MerchInput{
		Title: req.Title, Description: req.Description,
		Price: req.Price, ImageURL: req.ImageURL, BuyURL: req.BuyURL,
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

func (h *MerchHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req MerchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	item, errs, err := h.service.Update(id, services.MerchInput{
		Title: req.Title, Description: req.Description,
		Price: req.Price, ImageURL: req.ImageURL, BuyURL: req.BuyURL,
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

func (h *MerchHandler) Delete(c *gin.Context) {
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
