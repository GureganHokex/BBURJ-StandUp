package security

import "testing"

func TestSafeRedirectPath(t *testing.T) {
	tests := []struct {
		next, fallback, want string
	}{
		{"/admin", "/admin", "/admin"},
		{"/admin/events", "/admin", "/admin/events"},
		{"", "/admin", "/admin"},
		{"https://evil.com", "/admin", "/admin"},
		{"//evil.com", "/admin", "/admin"},
		{"/\\evil", "/admin", "/admin"},
		{"/admin?x=1", "/admin", "/admin?x=1"},
	}
	for _, tc := range tests {
		if got := SafeRedirectPath(tc.next, tc.fallback); got != tc.want {
			t.Errorf("SafeRedirectPath(%q, %q) = %q, want %q", tc.next, tc.fallback, got, tc.want)
		}
	}
}
