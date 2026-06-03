package config

import (
	"strings"
	"testing"
)

func TestValidateProduction(t *testing.T) {
	cfg := Config{
		AppEnv:        "production",
		SessionSecret: "short",
		AdminPassword: "admin123",
		SecureCookies: false,
		DatabaseURL:   "postgres://u:p@host/db?sslmode=disable",
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNormalizeDatabaseURLRender(t *testing.T) {
	in := "postgres://u:p@dpg-test-a.frankfurt-postgres.render.com/comic?sslmode=disable"
	got := NormalizeDatabaseURL(in)
	if strings.Contains(got, "sslmode=disable") {
		t.Fatalf("expected sslmode rewritten, got %q", got)
	}
	if !strings.Contains(got, "sslmode=prefer") {
		t.Fatalf("expected sslmode=prefer, got %q", got)
	}
}

func TestValidateDevelopmentSkips(t *testing.T) {
	cfg := Config{AppEnv: "development", AdminPassword: "admin123"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("dev should skip validation: %v", err)
	}
}
