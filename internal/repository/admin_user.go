package repository

import (
	"github.com/burj/comic/internal/models"
	"gorm.io/gorm"
)

type AdminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

func (r *AdminUserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.AdminUser{}).Count(&count).Error
	return count, err
}

func (r *AdminUserRepository) FindByUsername(username string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AdminUserRepository) Create(user *models.AdminUser) error {
	return r.db.Create(user).Error
}
