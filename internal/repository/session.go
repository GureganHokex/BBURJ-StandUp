package repository

import (
	"time"

	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(id string, adminID uint, expiresAt time.Time) error {
	return r.db.Create(&models.AdminSession{
		ID:        id,
		AdminID:   adminID,
		ExpiresAt: expiresAt,
	}).Error
}

func (r *SessionRepository) FindAdminID(id string) (uint, time.Time, error) {
	var session models.AdminSession
	err := r.db.Where("id = ?", id).First(&session).Error
	if err != nil {
		return 0, time.Time{}, err
	}
	return session.AdminID, session.ExpiresAt, nil
}

func (r *SessionRepository) Delete(id string) error {
	return r.db.Delete(&models.AdminSession{}, "id = ?", id).Error
}

func (r *SessionRepository) DeleteByAdminID(adminID uint) error {
	return r.db.Where("admin_id = ?", adminID).Delete(&models.AdminSession{}).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.AdminSession{}).Error
}
