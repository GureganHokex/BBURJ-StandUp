package repository

import (
	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type MerchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) *MerchRepository {
	return &MerchRepository{db: db}
}

func (r *MerchRepository) List(limit, offset int) ([]models.Merch, int64, error) {
	var total int64
	if err := r.db.Model(&models.Merch{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Merch
	err := r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *MerchRepository) FindByID(id uint) (*models.Merch, error) {
	var item models.Merch
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MerchRepository) Create(item *models.Merch) error {
	return r.db.Create(item).Error
}

func (r *MerchRepository) Update(item *models.Merch) error {
	return r.db.Save(item).Error
}

func (r *MerchRepository) Delete(id uint) error {
	return r.db.Delete(&models.Merch{}, id).Error
}

func (r *MerchRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Merch{}).Count(&count).Error
	return count, err
}
