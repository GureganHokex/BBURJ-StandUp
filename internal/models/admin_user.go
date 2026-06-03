package models

import "time"

type AdminUser struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:150;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
