package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultSessionSecret = "change-me-32-chars-minimum-secret"
	defaultAdminPassword = "admin123"
	minSessionSecretLen  = 32
	minAdminPasswordLen  = 12
)

type Config struct {
	AppEnv         string
	Port           string
	DatabaseURL    string
	SessionSecret  string
	AdminUsername  string
	AdminPassword  string
	SecureCookies  bool
	UploadDir      string
	MaxUploadMB    int
	TrustedProxies []string
}

func Load() Config {
	trusted := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	var proxies []string
	if trusted != "" {
		for _, p := range strings.Split(trusted, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				proxies = append(proxies, p)
			}
		}
	}

	return Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    NormalizeDatabaseURL(getEnv("DATABASE_URL", "postgres://comic:comic@localhost:5432/comic?sslmode=disable")),
		SessionSecret:  getEnv("SESSION_SECRET", defaultSessionSecret),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", defaultAdminPassword),
		SecureCookies:  getEnvBool("SECURE_COOKIES", false),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadMB:    getEnvInt("MAX_UPLOAD_MB", 10),
		TrustedProxies: proxies,
	}
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// NormalizeDatabaseURL adjusts Render internal Postgres URLs for production startup.
func NormalizeDatabaseURL(dsn string) string {
	if dsn == "" {
		return dsn
	}
	// Render private-network URLs often ship with sslmode=disable; TLS is not used there.
	if strings.Contains(dsn, "sslmode=disable") &&
		(strings.Contains(dsn, ".render.com") || strings.Contains(dsn, "dpg-")) {
		return strings.Replace(dsn, "sslmode=disable", "sslmode=prefer", 1)
	}
	return dsn
}

// Validate fails fast when production is misconfigured.
func (c Config) Validate() error {
	if !c.IsProduction() {
		return nil
	}

	var errs []string

	if len(c.SessionSecret) < minSessionSecretLen {
		errs = append(errs, fmt.Sprintf("SESSION_SECRET must be at least %d characters", minSessionSecretLen))
	}
	if c.SessionSecret == defaultSessionSecret || strings.Contains(c.SessionSecret, "change-me") {
		errs = append(errs, "SESSION_SECRET must be a unique random value in production")
	}

	if len(c.AdminPassword) < minAdminPasswordLen {
		errs = append(errs, fmt.Sprintf("ADMIN_PASSWORD must be at least %d characters for initial seed", minAdminPasswordLen))
	}
	if c.AdminPassword == defaultAdminPassword {
		errs = append(errs, "ADMIN_PASSWORD must not be the default in production")
	}

	if !c.SecureCookies {
		errs = append(errs, "SECURE_COOKIES must be true in production (HTTPS)")
	}

	if strings.Contains(c.DatabaseURL, "sslmode=disable") {
		errs = append(errs, "use sslmode=require (or verify-full) in DATABASE_URL for production")
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("production config: %s", strings.Join(errs, "; "))
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
