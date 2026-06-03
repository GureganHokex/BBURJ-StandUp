package config

import "testing"

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

func TestValidateDevelopmentSkips(t *testing.T) {
	cfg := Config{AppEnv: "development", AdminPassword: "admin123"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("dev should skip validation: %v", err)
	}
}
