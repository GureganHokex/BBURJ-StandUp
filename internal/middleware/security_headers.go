package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders sets baseline HTTP security headers for public site and admin.
func SecurityHeaders(production bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if production {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		// CSP: allow self, uploads, Google fonts, YouTube embeds, Tailwind CDN in admin
		c.Header("Content-Security-Policy", strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com https://unpkg.com",
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdn.tailwindcss.com",
			"font-src 'self' https://fonts.gstatic.com",
			"img-src 'self' data: https:",
			"frame-src 'self' https://www.youtube.com https://www.youtube-nocookie.com",
			"connect-src 'self'",
			"base-uri 'self'",
			"form-action 'self'",
			"frame-ancestors 'none'",
		}, "; "))
		c.Next()
	}
}
