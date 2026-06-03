package services

import (
	"errors"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"gorm.io/gorm"
)

type PhotoInput struct {
	Title     string
	ImageURL  string
	SortOrder int
}

type PhotoService struct {
	repo *repository.PhotoRepository
}

func NewPhotoService(repo *repository.PhotoRepository) *PhotoService {
	return &PhotoService{repo: repo}
}

func (s *PhotoService) List(limit, offset int) ([]models.Photo, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.List(limit, offset)
}

func (s *PhotoService) Get(id uint) (*models.Photo, error) {
	return s.repo.FindByID(id)
}

func (s *PhotoService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *PhotoService) Create(input PhotoInput) (*models.Photo, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}
	item := &models.Photo{
		Title: input.Title, ImageURL: input.ImageURL, SortOrder: input.SortOrder,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, nil, err
	}
	return item, nil, nil
}

func (s *PhotoService) Update(id uint, input PhotoInput) (*models.Photo, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}
	item, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		return nil, nil, err
	}
	item.Title = input.Title
	item.ImageURL = input.ImageURL
	item.SortOrder = input.SortOrder
	if err := s.repo.Update(item); err != nil {
		return nil, nil, err
	}
	return item, nil, nil
}

func (s *PhotoService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *PhotoService) validate(input PhotoInput) FieldErrors {
	return validateImageURL("image_url", input.ImageURL)
}
