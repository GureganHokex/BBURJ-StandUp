package api

import (
	"errors"
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	auth *services.AuthService
}

func NewAccountHandler(auth *services.AuthService) *AccountHandler {
	return &AccountHandler{auth: auth}
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *AccountHandler) ChangePassword(c *gin.Context) {
	adminID, ok := c.Get("admin_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}

	err := h.auth.ChangePassword(adminID.(uint), req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Errors: services.FieldErrors{
				"current_password": "incorrect password",
			}})
		case errors.Is(err, services.ErrWeakPassword):
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Errors: services.FieldErrors{
				"new_password": "must be at least 12 characters",
			}})
		default:
			writeInternalError(c, appConfig(c), err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
