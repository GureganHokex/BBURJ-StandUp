package services

import (
	"errors"
	"time"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"gorm.io/gorm"
)

type EventInput struct {
	Title       string
	Date        time.Time
	City        string
	Description string
	TicketURL   string
}

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) List(limit, offset int, upcomingOnly bool) ([]models.Event, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.List(limit, offset, upcomingOnly)
}

func (s *EventService) Get(id uint) (*models.Event, error) {
	return s.repo.FindByID(id)
}

func (s *EventService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *EventService) Create(input EventInput) (*models.Event, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}

	event := &models.Event{
		Title:       input.Title,
		Date:        input.Date,
		City:        input.City,
		Description: input.Description,
		TicketURL:   input.TicketURL,
	}
	if err := s.repo.Create(event); err != nil {
		return nil, nil, err
	}
	return event, nil, nil
}

func (s *EventService) Update(id uint, input EventInput) (*models.Event, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}

	event, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		return nil, nil, err
	}

	event.Title = input.Title
	event.Date = input.Date
	event.City = input.City
	event.Description = input.Description
	event.TicketURL = input.TicketURL

	if err := s.repo.Update(event); err != nil {
		return nil, nil, err
	}
	return event, nil, nil
}

func (s *EventService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *EventService) validate(input EventInput) FieldErrors {
	errs := validateRequired(map[string]string{
		"title": input.Title,
		"city":  input.City,
	})
	if input.Date.IsZero() {
		errs["date"] = "required"
	}
	return mergeErrors(errs, validateOptionalURL("ticket_url", input.TicketURL))
}
