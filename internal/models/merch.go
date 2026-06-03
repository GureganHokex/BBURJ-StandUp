package models

import "time"

func (Merch) TableName() string { return "merch" }

type Merch struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Price       int       `gorm:"not null" json:"price"`
	ImageURL    string    `gorm:"size:512" json:"image_url"`
	BuyURL      string    `gorm:"size:512" json:"buy_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
