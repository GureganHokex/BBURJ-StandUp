package api

import (
	"net/http"

	"github.com/burj/comic/internal/storage"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploader *storage.Uploader
}

func NewUploadHandler(uploader *storage.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "file is required"})
		return
	}

	url, err := h.uploader.Save(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ItemResponse[map[string]string]{
		Data: map[string]string{"url": url},
	})
}
