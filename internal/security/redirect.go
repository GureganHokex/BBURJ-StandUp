package security

import "strings"

// SafeRedirectPath allows only same-site relative paths (open redirect protection).
func SafeRedirectPath(next, fallback string) string {
	if fallback == "" {
		fallback = "/admin"
	}
	next = strings.TrimSpace(next)
	if next == "" {
		return fallback
	}
	if !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return fallback
	}
	if strings.Contains(next, "://") || strings.Contains(next, "\\") {
		return fallback
	}
	return next
}
