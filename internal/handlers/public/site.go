package public

import (
	"net/http"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

func mergeSite(settings *services.SiteSettingsService, data gin.H) gin.H {
	if data == nil {
		data = gin.H{}
	}
	s, err := settings.Get()
	if err == nil && s != nil {
		data["Settings"] = s
	} else {
		data["Settings"] = defaultSiteSettings()
	}
	return data
}

func defaultSiteSettings() *models.SiteSettings {
	return &models.SiteSettings{
		HeroTagline:     "СТЕНДАП-КОМИК ИЗ САНКТ-ПЕТЕРБУРГА",
		YouTubeURL:      "https://youtube.com",
		TelegramURL:     "https://t.me",
		InstagramURL:    "https://instagram.com",
		YouTubeHandle:   "@bburj",
		TelegramHandle:  "@bburj",
		InstagramHandle: "@bburj",
		ShowEvents:      true,
		ShowVideos:      true,
		ShowPhotos:      true,
		ShowMerch:       true,
		ShowAbout:       true,
	}
}

func sectionVisible(settings *services.SiteSettingsService, check func(*models.SiteSettings) bool) bool {
	s, err := settings.Get()
	if err != nil || s == nil {
		return true
	}
	return check(s)
}

func denyHiddenSection(c *gin.Context, settings *services.SiteSettingsService, check func(*models.SiteSettings) bool) bool {
	if sectionVisible(settings, check) {
		return true
	}
	c.Status(http.StatusNotFound)
	return false
}
