package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/burj/comic/internal/config"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func appConfig(c *gin.Context) config.Config {
	if v, ok := c.Get("cfg"); ok {
		if cfg, ok := v.(config.Config); ok {
			return cfg
		}
	}
	return config.Config{}
}

const maxPageLimit = 100

func parsePagination(c *gin.Context) (limit, offset int) {
	limit = 50
	offset = 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
			if limit > maxPageLimit {
				limit = maxPageLimit
			}
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return 0, false
	}
	return uint(id), true
}

func parseDate(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid date")
}

func writeValidationErrors(c *gin.Context, errs services.FieldErrors) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Errors: errs})
}

func writeNotFound(c *gin.Context, err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
		return true
	}
	return false
}
