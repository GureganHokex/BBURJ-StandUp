package public

import (
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
		data["Settings"] = &models.SiteSettings{
			HeroTagline:     "СТЕНДАП-КОМИК ИЗ САНКТ-ПЕТЕРБУРГА",
			YouTubeURL:      "https://youtube.com",
			TelegramURL:     "https://t.me",
			InstagramURL:    "https://instagram.com",
			YouTubeHandle:   "@bburj",
			TelegramHandle:  "@bburj",
			InstagramHandle: "@bburj",
		}
	}
	return data
}
