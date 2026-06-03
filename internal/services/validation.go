package services

import (
	"net/url"
	"strings"
)

type FieldErrors map[string]string

func (e FieldErrors) HasErrors() bool {
	return len(e) > 0
}

func validateRequired(fields map[string]string) FieldErrors {
	errs := FieldErrors{}
	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			errs[field] = "required"
		}
	}
	return errs
}

func validateOptionalURL(field, value string) FieldErrors {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if isLocalAssetPath(value) {
		return nil
	}
	if _, err := url.ParseRequestURI(value); err != nil {
		return FieldErrors{field: "invalid url"}
	}
	return nil
}

func validateImageURL(field, value string) FieldErrors {
	if strings.TrimSpace(value) == "" {
		return FieldErrors{field: "required"}
	}
	return validateOptionalURL(field, value)
}

func isLocalAssetPath(value string) bool {
	if !strings.HasPrefix(value, "/") {
		return false
	}
	if strings.Contains(value, "..") {
		return false
	}
	return strings.HasPrefix(value, "/uploads/") || strings.HasPrefix(value, "/static/")
}

func mergeErrors(errs ...FieldErrors) FieldErrors {
	out := FieldErrors{}
	for _, e := range errs {
		for k, v := range e {
			out[k] = v
		}
	}
	return out
}
