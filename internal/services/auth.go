package services

import (
	"errors"

	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/repository"
	"github.com/burj/comic/internal/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	adminRepo *repository.AdminUserRepository
	sessions  *session.Store
}

func NewAuthService(adminRepo *repository.AdminUserRepository, sessions *session.Store) *AuthService {
	return &AuthService{adminRepo: adminRepo, sessions: sessions}
}

func (s *AuthService) SeedAdmin(username, password string) error {
	count, err := s.adminRepo.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.adminRepo.Create(&models.AdminUser{
		Username:     username,
		PasswordHash: string(hash),
	})
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.sessions.Create(user.ID), nil
}

func (s *AuthService) Logout(sessionID string) {
	s.sessions.Delete(sessionID)
}

func (s *AuthService) GetAdminID(sessionID string) (uint, bool) {
	return s.sessions.GetAdminID(sessionID)
}
