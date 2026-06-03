package middleware

import (
	"net/http"

	"github.com/burj/comic/internal/config"
	"github.com/gin-gonic/gin"
)

const SessionCookieName = "comic_session"

func SetSessionCookie(c *gin.Context, cfg config.Config, sessionID string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SessionCookieName, sessionID, 86400, "/", "", cfg.SecureCookies, true)
}

func ClearSessionCookie(c *gin.Context, cfg config.Config) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SessionCookieName, "", -1, "/", "", cfg.SecureCookies, true)
}

func SessionID(c *gin.Context) string {
	id, err := c.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}
	return id
}
