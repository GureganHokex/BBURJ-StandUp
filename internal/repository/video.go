package repository

import (
	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) List(limit, offset int) ([]models.Video, int64, error) {
	var total int64
	if err := r.db.Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var videos []models.Video
	err := r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&videos).Error
	return videos, total, err
}

func (r *VideoRepository) FindByID(id uint) (*models.Video, error) {
	var video models.Video
	err := r.db.First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *VideoRepository) Create(video *models.Video) error {
	return r.db.Create(video).Error
}

func (r *VideoRepository) Update(video *models.Video) error {
	return r.db.Save(video).Error
}

func (r *VideoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Video{}, id).Error
}

func (r *VideoRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Video{}).Count(&count).Error
	return count, err
}
