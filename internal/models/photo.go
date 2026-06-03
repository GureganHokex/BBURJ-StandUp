package models

import "time"

type Photo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255" json:"title"`
	ImageURL  string    `gorm:"size:512;not null" json:"image_url"`
	SortOrder int       `gorm:"default:0;index" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
