package main

import (
	"log"
	"os"

	"github.com/burj/comic/internal/config"
	"github.com/burj/comic/internal/database"
	"github.com/joho/godotenv"
)

// Run SQL migrations manually: go run ./cmd/migrate
func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL, cfg.AppEnv)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	dir := "migrations"
	if v := os.Getenv("MIGRATIONS_DIR"); v != "" {
		dir = v
	}

	if err := database.RunSQLMigrations(db, dir); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied successfully")
}
