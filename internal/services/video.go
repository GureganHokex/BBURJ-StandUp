package services

import (
	"errors"
	"net/url"
	"strings"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"gorm.io/gorm"
)

type VideoInput struct {
	Title string
	URL   string
}

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService(repo *repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) List(limit, offset int) ([]models.Video, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.List(limit, offset)
}

func (s *VideoService) Get(id uint) (*models.Video, error) {
	return s.repo.FindByID(id)
}

func (s *VideoService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *VideoService) Create(input VideoInput) (*models.Video, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}

	video := &models.Video{
		Title:    input.Title,
		URL:      input.URL,
		Platform: DetectPlatform(input.URL),
	}
	if err := s.repo.Create(video); err != nil {
		return nil, nil, err
	}
	return video, nil, nil
}

func (s *VideoService) Update(id uint, input VideoInput) (*models.Video, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}

	video, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		return nil, nil, err
	}

	video.Title = input.Title
	video.URL = input.URL
	video.Platform = DetectPlatform(input.URL)

	if err := s.repo.Update(video); err != nil {
		return nil, nil, err
	}
	return video, nil, nil
}

func (s *VideoService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *VideoService) validate(input VideoInput) FieldErrors {
	errs := validateRequired(map[string]string{
		"title": input.Title,
		"url":   input.URL,
	})
	if _, err := url.ParseRequestURI(input.URL); input.URL != "" && err != nil {
		errs["url"] = "invalid url"
	}
	return errs
}

func DetectPlatform(rawURL string) models.VideoPlatform {
	lower := strings.ToLower(rawURL)
	switch {
	case strings.Contains(lower, "youtube.com"), strings.Contains(lower, "youtu.be"):
		return models.PlatformYouTube
	case strings.Contains(lower, "rutube.ru"):
		return models.PlatformRuTube
	case strings.Contains(lower, "vk.com"), strings.Contains(lower, "vkvideo.ru"):
		return models.PlatformVK
	default:
		return models.PlatformOther
	}
}

func EmbedURL(platform models.VideoPlatform, rawURL string) string {
	lower := strings.ToLower(rawURL)
	switch platform {
	case models.PlatformYouTube:
		if strings.Contains(lower, "youtu.be/") {
			parts := strings.Split(lower, "youtu.be/")
			if len(parts) > 1 {
				id := strings.Split(parts[1], "?")[0]
				return "https://www.youtube.com/embed/" + id
			}
		}
		if strings.Contains(lower, "v=") {
			parts := strings.Split(lower, "v=")
			if len(parts) > 1 {
				id := strings.Split(parts[1], "&")[0]
				return "https://www.youtube.com/embed/" + id
			}
		}
	case models.PlatformRuTube:
		if strings.Contains(lower, "/video/") {
			parts := strings.Split(lower, "/video/")
			if len(parts) > 1 {
				id := strings.Split(parts[1], "/")[0]
				return "https://rutube.ru/play/embed/" + id
			}
		}
	case models.PlatformVK:
		return rawURL
	}
	return rawURL
}
