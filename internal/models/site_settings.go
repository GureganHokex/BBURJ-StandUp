package models

import "time"

type SiteSettings struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	HeroImageURL      string    `gorm:"size:512" json:"hero_image_url"`
	PortraitImageURL  string    `gorm:"size:512" json:"portrait_image_url"`
	HeroTagline       string    `gorm:"size:255" json:"hero_tagline"`
	AboutText         string    `gorm:"type:text" json:"about_text"`
	AboutExtra        string    `gorm:"type:text" json:"about_extra"`
	YouTubeURL        string    `gorm:"column:youtube_url;size:512" json:"youtube_url"`
	TelegramURL       string    `gorm:"size:512" json:"telegram_url"`
	InstagramURL      string    `gorm:"size:512" json:"instagram_url"`
	YouTubeHandle     string    `gorm:"column:youtube_handle;size:128" json:"youtube_handle"`
	TelegramHandle    string    `gorm:"size:128" json:"telegram_handle"`
	InstagramHandle   string    `gorm:"size:128" json:"instagram_handle"`
	ContactEmail      string    `gorm:"size:255" json:"contact_email"`
	ContactPhone      string    `gorm:"size:64" json:"contact_phone"`
	ContactTelegram   string    `gorm:"size:128" json:"contact_telegram"`
	TimepadOrgID          string `gorm:"column:timepad_org_id;size:32" json:"timepad_org_id"`
	TimepadAPIKey         string `gorm:"column:timepad_api_key;size:255" json:"-"`
	TicketscloudOrgID     string `gorm:"column:ticketscloud_org_id;size:128" json:"ticketscloud_org_id"`
	TicketscloudAPIKey    string `gorm:"column:ticketscloud_api_key;size:255" json:"-"`
	EventImportKeywords   string `gorm:"column:event_import_keywords;size:512" json:"event_import_keywords"`
	ShowEvents            bool   `gorm:"not null;default:true" json:"show_events"`
	ShowVideos            bool   `gorm:"not null;default:true" json:"show_videos"`
	ShowPhotos            bool   `gorm:"not null;default:true" json:"show_photos"`
	ShowMerch             bool   `gorm:"not null;default:true" json:"show_merch"`
	ShowAbout             bool   `gorm:"not null;default:true" json:"show_about"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
