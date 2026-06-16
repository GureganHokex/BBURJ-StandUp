package api

import (
	"log"
	"net/http"

	"github.com/burj/comic/internal/config"
	"github.com/gin-gonic/gin"
)

func writeInternalError(c *gin.Context, cfg config.Config, err error) {
	if err != nil {
		log.Printf("api error %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	}
	if cfg.IsProduction() {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
}
