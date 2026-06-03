package services

import (
	"errors"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"gorm.io/gorm"
)

type MerchInput struct {
	Title       string
	Description string
	Price       int
	ImageURL    string
	BuyURL      string
}

type MerchService struct {
	repo *repository.MerchRepository
}

func NewMerchService(repo *repository.MerchRepository) *MerchService {
	return &MerchService{repo: repo}
}

func (s *MerchService) List(limit, offset int) ([]models.Merch, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.List(limit, offset)
}

func (s *MerchService) Get(id uint) (*models.Merch, error) {
	return s.repo.FindByID(id)
}

func (s *MerchService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *MerchService) Create(input MerchInput) (*models.Merch, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}

	item := &models.Merch{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		ImageURL:    input.ImageURL,
		BuyURL:      input.BuyURL,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, nil, err
	}
	return item, nil, nil
}

func (s *MerchService) Update(id uint, input MerchInput) (*models.Merch, FieldErrors, error) {
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
	item.Description = input.Description
	item.Price = input.Price
	item.ImageURL = input.ImageURL
	item.BuyURL = input.BuyURL

	if err := s.repo.Update(item); err != nil {
		return nil, nil, err
	}
	return item, nil, nil
}

func (s *MerchService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *MerchService) validate(input MerchInput) FieldErrors {
	errs := validateRequired(map[string]string{
		"title": input.Title,
	})
	if input.Price <= 0 {
		errs["price"] = "must be greater than 0"
	}
	imgErrs := FieldErrors{}
	if input.ImageURL != "" {
		imgErrs = validateOptionalURL("image_url", input.ImageURL)
	}
	return mergeErrors(errs, imgErrs, validateOptionalURL("buy_url", input.BuyURL))
}
