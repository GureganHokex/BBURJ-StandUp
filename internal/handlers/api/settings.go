package api

import (
	"net/http"

	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	service *services.SiteSettingsService
}

func NewSettingsHandler(service *services.SiteSettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

type SettingsRequest struct {
	HeroImageURL     string `json:"hero_image_url"`
	PortraitImageURL string `json:"portrait_image_url"`
	HeroTagline      string `json:"hero_tagline"`
	AboutText        string `json:"about_text"`
	AboutExtra       string `json:"about_extra"`
	YouTubeURL       string `json:"youtube_url"`
	TelegramURL      string `json:"telegram_url"`
	InstagramURL     string `json:"instagram_url"`
	YouTubeHandle    string `json:"youtube_handle"`
	TelegramHandle   string `json:"telegram_handle"`
	InstagramHandle  string `json:"instagram_handle"`
	ContactEmail     string `json:"contact_email"`
	ContactPhone     string `json:"contact_phone"`
	ContactTelegram  string `json:"contact_telegram"`
	TimepadOrgID          string `json:"timepad_org_id"`
	TimepadAPIKey         string `json:"timepad_api_key"`
	TicketscloudOrgID     string `json:"ticketscloud_org_id"`
	TicketscloudAPIKey    string `json:"ticketscloud_api_key"`
	EventImportKeywords   string `json:"event_import_keywords"`
	ShowEvents            bool   `json:"show_events"`
	ShowVideos            bool   `json:"show_videos"`
	ShowPhotos            bool   `json:"show_photos"`
	ShowMerch             bool   `json:"show_merch"`
	ShowAbout             bool   `json:"show_about"`
}

func (h *SettingsHandler) Get(c *gin.Context) {
	s, err := h.service.Get()
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: s})
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var req SettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}
	s, errs, err := h.service.Update(services.SiteSettingsInput{
		HeroImageURL: req.HeroImageURL, PortraitImageURL: req.PortraitImageURL,
		HeroTagline: req.HeroTagline, AboutText: req.AboutText, AboutExtra: req.AboutExtra,
		YouTubeURL: req.YouTubeURL, TelegramURL: req.TelegramURL, InstagramURL: req.InstagramURL,
		YouTubeHandle: req.YouTubeHandle, TelegramHandle: req.TelegramHandle,
		InstagramHandle: req.InstagramHandle,
		ContactEmail: req.ContactEmail, ContactPhone: req.ContactPhone,
		ContactTelegram: req.ContactTelegram,
		TimepadOrgID: req.TimepadOrgID, TimepadAPIKey: req.TimepadAPIKey,
		TicketscloudOrgID: req.TicketscloudOrgID, TicketscloudAPIKey: req.TicketscloudAPIKey,
		EventImportKeywords: req.EventImportKeywords,
		ShowEvents: req.ShowEvents, ShowVideos: req.ShowVideos,
		ShowPhotos: req.ShowPhotos, ShowMerch: req.ShowMerch, ShowAbout: req.ShowAbout,
	})
	if err != nil {
		writeInternalError(c, appConfig(c), err)
		return
	}
	if errs != nil && errs.HasErrors() {
		writeValidationErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, ItemResponse[any]{Data: s})
}
