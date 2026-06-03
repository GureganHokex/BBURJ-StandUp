package repository

import (
	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type PhotoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (r *PhotoRepository) List(limit, offset int) ([]models.Photo, int64, error) {
	var total int64
	if err := r.db.Model(&models.Photo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.Photo
	err := r.db.Order("sort_order ASC, id DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *PhotoRepository) FindByID(id uint) (*models.Photo, error) {
	var item models.Photo
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PhotoRepository) Create(item *models.Photo) error {
	return r.db.Create(item).Error
}

func (r *PhotoRepository) Update(item *models.Photo) error {
	return r.db.Save(item).Error
}

func (r *PhotoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Photo{}, id).Error
}

func (r *PhotoRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Photo{}).Count(&count).Error
	return count, err
}
