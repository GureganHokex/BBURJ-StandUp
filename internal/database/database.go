package database

import (
	"fmt"

	"github.com/burj/comic/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string, appEnv string) (*gorm.DB, error) {
	logLevel := logger.Info
	if appEnv == "production" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.AdminUser{},
		&models.AdminSession{},
		&models.Event{},
		&models.Video{},
		&models.Merch{},
		&models.Photo{},
		&models.SiteSettings{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
