package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is a simple in-memory per-key sliding window limiter.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string][]time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{entries: make(map[string][]time.Time)}
}

func (l *RateLimiter) Allow(key string, max int, window time.Duration) bool {
	if max <= 0 || window <= 0 {
		return true
	}
	now := time.Now()
	cutoff := now.Add(-window)

	l.mu.Lock()
	defer l.mu.Unlock()

	times := l.entries[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= max {
		l.entries[key] = filtered
		return false
	}
	filtered = append(filtered, now)
	l.entries[key] = filtered
	return true
}

func RateLimit(l *RateLimiter, max int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		if !l.Allow(key, max, window) {
			if isAPI(c) {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
				return
			}
			c.String(http.StatusTooManyRequests, "Too many requests. Try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isAPI(c *gin.Context) bool {
	return len(c.Request.URL.Path) >= 5 && c.Request.URL.Path[:5] == "/api/"
}
