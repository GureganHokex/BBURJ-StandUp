package services

import (
	"errors"
	"strings"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"gorm.io/gorm"
)

type SiteSettingsInput struct {
	HeroImageURL     string
	PortraitImageURL string
	HeroTagline      string
	AboutText        string
	AboutExtra       string
	YouTubeURL       string
	TelegramURL      string
	InstagramURL     string
	YouTubeHandle    string
	TelegramHandle   string
	InstagramHandle  string
	TimepadOrgID          string
	TimepadAPIKey         string
	TicketscloudOrgID     string
	TicketscloudAPIKey    string
	EventImportKeywords   string
	ShowEvents            bool
	ShowVideos            bool
	ShowPhotos            bool
	ShowMerch             bool
	ShowAbout             bool
}

type SiteSettingsService struct {
	repo *repository.SiteSettingsRepository
}

func NewSiteSettingsService(repo *repository.SiteSettingsRepository) *SiteSettingsService {
	return &SiteSettingsService{repo: repo}
}

func (s *SiteSettingsService) SeedDefaults() (*models.SiteSettings, error) {
	return s.repo.EnsureDefault(models.SiteSettings{
		ID:              1,
		HeroImageURL:    "/static/img/hero.jpg",
		PortraitImageURL: "/static/img/portrait.jpg",
		HeroTagline:     "СТЕНДАП-КОМИК ИЗ САНКТ-ПЕТЕРБУРГА",
		AboutText:       "Большой Буржинский — стендап-комик из Санкт-Петербурга. Выступает с авторским материалом: наблюдения, истории из жизни и всё, что бесит — но смешно.",
		AboutExtra:      "Приходи на живые шоу или смотри записи — тут собраны афиша, видео и мерч.",
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
	})
}

func (s *SiteSettingsService) Get() (*models.SiteSettings, error) {
	settings, err := s.repo.Get()
	if err == nil {
		return settings, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.SeedDefaults()
	}
	return nil, err
}

func (s *SiteSettingsService) Update(input SiteSettingsInput) (*models.SiteSettings, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}
	current, err := s.Get()
	if err != nil {
		return nil, nil, err
	}
	current.HeroImageURL = input.HeroImageURL
	current.PortraitImageURL = input.PortraitImageURL
	current.HeroTagline = input.HeroTagline
	current.AboutText = input.AboutText
	current.AboutExtra = input.AboutExtra
	current.YouTubeURL = input.YouTubeURL
	current.TelegramURL = input.TelegramURL
	current.InstagramURL = input.InstagramURL
	current.YouTubeHandle = input.YouTubeHandle
	current.TelegramHandle = input.TelegramHandle
	current.InstagramHandle = input.InstagramHandle
	current.TimepadOrgID = input.TimepadOrgID
	current.TicketscloudOrgID = input.TicketscloudOrgID
	current.EventImportKeywords = input.EventImportKeywords
	current.ShowEvents = input.ShowEvents
	current.ShowVideos = input.ShowVideos
	current.ShowPhotos = input.ShowPhotos
	current.ShowMerch = input.ShowMerch
	current.ShowAbout = input.ShowAbout
	if strings.TrimSpace(input.TimepadAPIKey) != "" {
		current.TimepadAPIKey = input.TimepadAPIKey
	}
	if strings.TrimSpace(input.TicketscloudAPIKey) != "" {
		current.TicketscloudAPIKey = input.TicketscloudAPIKey
	}
	if err := s.repo.Save(current); err != nil {
		return nil, nil, err
	}
	return current, nil, nil
}

func (s *SiteSettingsService) validate(input SiteSettingsInput) FieldErrors {
	errs := validateRequired(map[string]string{
		"hero_tagline": input.HeroTagline,
		"about_text":   input.AboutText,
	})
	return mergeErrors(
		errs,
		validateImageURL("hero_image_url", input.HeroImageURL),
		validateImageURL("portrait_image_url", input.PortraitImageURL),
		validateOptionalURL("youtube_url", input.YouTubeURL),
		validateOptionalURL("telegram_url", input.TelegramURL),
		validateOptionalURL("instagram_url", input.InstagramURL),
	)
}

var ErrSettingsNotFound = errors.New("settings not found")
