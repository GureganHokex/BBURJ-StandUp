package models

import "time"

type VideoPlatform string

const (
	PlatformYouTube VideoPlatform = "youtube"
	PlatformRuTube  VideoPlatform = "rutube"
	PlatformVK      VideoPlatform = "vk"
	PlatformOther   VideoPlatform = "other"
)

type Video struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Title     string        `gorm:"size:255;not null" json:"title"`
	URL       string        `gorm:"size:512;not null" json:"url"`
	Platform  VideoPlatform `gorm:"size:32;not null" json:"platform"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
