package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"github.com/burj/comic/internal/storage"
	"gorm.io/gorm"
)

var ErrDuplicateExternalEvent = errors.New("event already imported")

type EventInput struct {
	Title        string
	Date         time.Time
	City         string
	Description    string
	TicketURL      string
	PosterImageURL string
	TicketSource   string
	ExternalID   string
}

type EventService struct {
	repo     *repository.EventRepository
	uploader *storage.Uploader
}

func NewEventService(repo *repository.EventRepository, uploader *storage.Uploader) *EventService {
	return &EventService{repo: repo, uploader: uploader}
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

func (s *EventService) ExternalIDs(source string) (map[string]struct{}, error) {
	return s.repo.ExternalIDs(source)
}

func (s *EventService) Create(input EventInput) (*models.Event, FieldErrors, error) {
	if errs := s.validate(input); errs.HasErrors() {
		return nil, errs, nil
	}
	if err := s.ensureNotDuplicate(input); err != nil {
		return nil, nil, err
	}

	source := normalizeTicketSource(input.TicketSource)
	posterURL, err := s.mirrorPoster(context.Background(), input.PosterImageURL)
	if err != nil {
		return nil, nil, err
	}
	event := &models.Event{
		Title:          input.Title,
		Date:           input.Date,
		City:           input.City,
		Description:    input.Description,
		TicketURL:      input.TicketURL,
		PosterImageURL: posterURL,
		TicketSource:   source,
		ExternalID:     stringsTrimExternal(input.ExternalID, source),
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

	source := normalizeTicketSource(input.TicketSource)
	externalID := stringsTrimExternal(input.ExternalID, source)
	if source != event.TicketSource || externalID != event.ExternalID {
		if err := s.ensureNotDuplicateExcept(input, id); err != nil {
			return nil, nil, err
		}
	}

	event.Title = input.Title
	event.Date = input.Date
	event.City = input.City
	event.Description = input.Description
	event.TicketURL = input.TicketURL
	posterURL, err := s.mirrorPoster(context.Background(), input.PosterImageURL)
	if err != nil {
		return nil, nil, err
	}
	event.PosterImageURL = posterURL
	event.TicketSource = source
	event.ExternalID = externalID

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
	return mergeErrors(errs, validateOptionalURL("ticket_url", input.TicketURL), validateOptionalURL("poster_image_url", input.PosterImageURL))
}

func (s *EventService) ensureNotDuplicate(input EventInput) error {
	source := normalizeTicketSource(input.TicketSource)
	externalID := stringsTrimExternal(input.ExternalID, source)
	if externalID == "" {
		return nil
	}
	_, err := s.repo.FindByExternal(source, externalID)
	if err == nil {
		return ErrDuplicateExternalEvent
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (s *EventService) ensureNotDuplicateExcept(input EventInput, id uint) error {
	source := normalizeTicketSource(input.TicketSource)
	externalID := stringsTrimExternal(input.ExternalID, source)
	if externalID == "" {
		return nil
	}
	existing, err := s.repo.FindByExternal(source, externalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if existing.ID != id {
		return ErrDuplicateExternalEvent
	}
	return nil
}

func normalizeTicketSource(source string) string {
	if source == "" {
		return "manual"
	}
	return source
}

func stringsTrimExternal(id, source string) string {
	if source == "" || source == "manual" {
		return ""
	}
	return strings.TrimSpace(id)
}

func NormalizeEventForDisplay(e *models.Event) {
	parsed := ParseEventTitle(e.Title, e.Description)
	e.Title = parsed.Title
	e.Description = parsed.Description
}

func NormalizeEventsForDisplay(events []models.Event) []models.Event {
	for i := range events {
		NormalizeEventForDisplay(&events[i])
	}
	return events
}

func (s *EventService) mirrorPoster(ctx context.Context, posterURL string) (string, error) {
	if s.uploader == nil {
		return strings.TrimSpace(posterURL), nil
	}
	return s.uploader.ImportFromURL(ctx, posterURL)
}
