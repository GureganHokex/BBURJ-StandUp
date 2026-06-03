package models

import "time"

type SiteSettings struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	HeroImageURL      string    `gorm:"size:512" json:"hero_image_url"`
	PortraitImageURL  string    `gorm:"size:512" json:"portrait_image_url"`
	HeroTagline       string    `gorm:"size:255" json:"hero_tagline"`
	AboutText         string    `gorm:"type:text" json:"about_text"`
	AboutExtra        string    `gorm:"type:text" json:"about_extra"`
	YouTubeURL        string    `gorm:"size:512" json:"youtube_url"`
	TelegramURL       string    `gorm:"size:512" json:"telegram_url"`
	InstagramURL      string    `gorm:"size:512" json:"instagram_url"`
	YouTubeHandle     string    `gorm:"size:128" json:"youtube_handle"`
	TelegramHandle    string    `gorm:"size:128" json:"telegram_handle"`
	InstagramHandle   string    `gorm:"size:128" json:"instagram_handle"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
