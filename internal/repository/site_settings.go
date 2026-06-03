package repository

import (
	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type SiteSettingsRepository struct {
	db *gorm.DB
}

func NewSiteSettingsRepository(db *gorm.DB) *SiteSettingsRepository {
	return &SiteSettingsRepository{db: db}
}

func (r *SiteSettingsRepository) Get() (*models.SiteSettings, error) {
	var s models.SiteSettings
	err := r.db.First(&s, 1).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SiteSettingsRepository) Save(s *models.SiteSettings) error {
	s.ID = 1
	return r.db.Save(s).Error
}

func (r *SiteSettingsRepository) EnsureDefault(def models.SiteSettings) (*models.SiteSettings, error) {
	var s models.SiteSettings
	err := r.db.First(&s, 1).Error
	if err == nil {
		return &s, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	def.ID = 1
	if err := r.db.Create(&def).Error; err != nil {
		return nil, err
	}
	return &def, nil
}
