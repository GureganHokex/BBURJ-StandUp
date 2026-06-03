package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/burj/comic/internal/config"
	"github.com/gin-gonic/gin"
)

const CSRFCookieName = "csrf_token"
const CSRFFormField = "csrf_token"
const CSRFHeaderName = "X-CSRF-Token"

type CSRF struct {
	secret string
	secure bool
}

func NewCSRF(cfg config.Config) *CSRF {
	return &CSRF{secret: cfg.SessionSecret, secure: cfg.SecureCookies}
}

func (m *CSRF) EnsureToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Cookie(CSRFCookieName); err != nil {
			token := m.generateToken()
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie(CSRFCookieName, token, 86400, "/", "", m.secure, false)
		}
		c.Next()
	}
}

func (m *CSRF) Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		cookieToken, _ := c.Cookie(CSRFCookieName)
		formToken := c.PostForm(CSRFFormField)
		headerToken := c.GetHeader(CSRFHeaderName)
		submitted := formToken
		if submitted == "" {
			submitted = headerToken
		}

		if cookieToken == "" || submitted == "" || !m.valid(cookieToken, submitted) {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "csrf validation failed"})
				return
			}
			c.String(http.StatusForbidden, "CSRF validation failed")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *CSRF) TokenFromContext(c *gin.Context) string {
	token, err := c.Cookie(CSRFCookieName)
	if err != nil || token == "" {
		token = m.generateToken()
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(CSRFCookieName, token, 86400, "/", "", m.secure, false)
	}
	return token
}

func (m *CSRF) generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	raw := base64.RawURLEncoding.EncodeToString(b)
	mac := hmac.New(sha256.New, []byte(m.secret))
	_, _ = mac.Write([]byte(raw))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return raw + "." + sig
}

func (m *CSRF) valid(cookieToken, submitted string) bool {
	return hmac.Equal([]byte(cookieToken), []byte(submitted))
}
