package services

import "testing"

func TestValidateOptionalURLSchemes(t *testing.T) {
	cases := []struct {
		value string
		ok    bool
	}{
		{"", true},
		{"/uploads/abc.jpg", true},
		{"https://example.com", true},
		{"http://example.com/path", true},
		{"javascript:alert(1)", false},
		{"data:text/html,hi", false},
		{"ftp://example.com", false},
	}

	for _, tc := range cases {
		errs := validateOptionalURL("url", tc.value)
		hasErr := errs != nil && errs.HasErrors()
		if hasErr == tc.ok {
			t.Errorf("validateOptionalURL(%q): hasErr=%v, want %v", tc.value, hasErr, !tc.ok)
		}
	}
}
