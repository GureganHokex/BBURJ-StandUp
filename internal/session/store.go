package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/burj/comic/internal/repository"
	"gorm.io/gorm"
)

type Store struct {
	repo *repository.SessionRepository
	ttl  time.Duration
}

func NewStore(repo *repository.SessionRepository) *Store {
	return &Store{
		repo: repo,
		ttl:  24 * time.Hour,
	}
}

func (s *Store) Create(adminID uint) (string, error) {
	id := randomID()
	expiresAt := time.Now().Add(s.ttl)
	if err := s.repo.Create(id, adminID, expiresAt); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Store) GetAdminID(sessionID string) (uint, bool) {
	if sessionID == "" {
		return 0, false
	}
	adminID, expiresAt, err := s.repo.FindAdminID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, false
		}
		return 0, false
	}
	if time.Now().After(expiresAt) {
		_ = s.repo.Delete(sessionID)
		return 0, false
	}
	return adminID, true
}

func (s *Store) Delete(sessionID string) {
	if sessionID == "" {
		return
	}
	_ = s.repo.Delete(sessionID)
}

func (s *Store) DeleteAllForAdmin(adminID uint) {
	_ = s.repo.DeleteByAdminID(adminID)
}

func (s *Store) CleanupExpired() {
	_ = s.repo.DeleteExpired()
}

func randomID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
