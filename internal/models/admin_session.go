package models

import "time"

type AdminSession struct {
	ID        string    `gorm:"primaryKey;size:64"`
	AdminID   uint      `gorm:"not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}
