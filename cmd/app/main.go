package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/burj/comic/internal/config"
	"github.com/burj/comic/internal/database"
	"github.com/burj/comic/internal/handlers"
	"github.com/burj/comic/internal/repository"
	"github.com/burj/comic/internal/services"
	"github.com/burj/comic/internal/session"
	"github.com/burj/comic/internal/storage"
	"github.com/burj/comic/web"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config: %v", err)
	}

	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = database.Connect(cfg.DatabaseURL, cfg.AppEnv)
		if err == nil {
			break
		}
		log.Printf("database connect attempt %d/10: %v", attempt, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	migrationsDir := getEnv("MIGRATIONS_DIR", "migrations")
	if cfg.IsProduction() {
		if err := database.RunSQLMigrations(db, migrationsDir); err != nil {
			log.Fatalf("migrations: %v", err)
		}
	} else if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	adminRepo := repository.NewAdminUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	videoRepo := repository.NewVideoRepository(db)
	merchRepo := repository.NewMerchRepository(db)
	photoRepo := repository.NewPhotoRepository(db)
	settingsRepo := repository.NewSiteSettingsRepository(db)

	sessionRepo := repository.NewSessionRepository(db)
	sessions := session.NewStore(sessionRepo)
	sessions.CleanupExpired()
	authService := services.NewAuthService(adminRepo, sessions)

	uploader, err := storage.NewUploader(cfg.UploadDir, cfg.MaxUploadMB)
	if err != nil {
		log.Fatalf("uploader: %v", err)
	}

	eventService := services.NewEventService(eventRepo, uploader)
	videoService := services.NewVideoService(videoRepo)
	merchService := services.NewMerchService(merchRepo)
	photoService := services.NewPhotoService(photoRepo)
	settingsService := services.NewSiteSettingsService(settingsRepo)
	urlPreviewService := services.NewURLPreviewService()

	if err := authService.SeedAdmin(cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Fatalf("seed admin: %v", err)
	}
	if _, err := settingsService.SeedDefaults(); err != nil {
		log.Fatalf("seed settings: %v", err)
	}

	router := handlers.NewRouter(handlers.Deps{
		Config:    cfg,
		Auth:      authService,
		Events:    eventService,
		Videos:    videoService,
		Merch:     merchService,
		Photos:    photoService,
		Settings:   settingsService,
		URLPreview: urlPreviewService,
		Uploader:  uploader,
		UploadDir: cfg.UploadDir,
		Templates: web.Templates(),
		StaticFS:  web.StaticFS(),
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
