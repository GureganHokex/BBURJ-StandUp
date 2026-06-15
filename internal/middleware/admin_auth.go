package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

func AdminAuth(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if isLoginPath(path) {
			c.Next()
			return
		}

		sessionID := SessionID(c)
		adminID, ok := auth.GetAdminID(sessionID)
		if !ok {
			if strings.HasPrefix(path, "/api/") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			next := path
			if c.Request.URL.RawQuery != "" {
				next += "?" + c.Request.URL.RawQuery
			}
			c.Redirect(http.StatusFound, "/admin/login?next="+url.QueryEscape(next))
			c.Abort()
			return
		}

		c.Set("admin_id", adminID)
		c.Next()
	}
}

func isLoginPath(path string) bool {
	return path == "/admin/login"
}
