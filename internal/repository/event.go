package repository

import (
	"time"

	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) List(limit, offset int, upcomingOnly bool) ([]models.Event, int64, error) {
	q := r.db.Model(&models.Event{})
	if upcomingOnly {
		q = q.Where("date >= ?", time.Now())
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var events []models.Event
	err := q.Order("date ASC").Limit(limit).Offset(offset).Find(&events).Error
	return events, total, err
}

func (r *EventRepository) FindByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *EventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}

func (r *EventRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Event{}).Count(&count).Error
	return count, err
}
