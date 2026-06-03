package models

import "time"

type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Date        time.Time `gorm:"not null;index" json:"date"`
	City        string    `gorm:"size:128;not null" json:"city"`
	Description string    `gorm:"type:text" json:"description"`
	TicketURL   string    `gorm:"size:512" json:"ticket_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
