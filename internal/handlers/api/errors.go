package api

import (
	"net/http"

	"github.com/burj/comic/internal/config"
	"github.com/gin-gonic/gin"
)

func writeInternalError(c *gin.Context, cfg config.Config, err error) {
	if cfg.IsProduction() {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
}
