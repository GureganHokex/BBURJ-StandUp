package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type data struct {
	AdminID   uint
	ExpiresAt time.Time
}

type Store struct {
	secret string
	ttl    time.Duration
	mu     sync.RWMutex
	items  map[string]data
}

func NewStore(secret string) *Store {
	return &Store{
		secret: secret,
		ttl:    24 * time.Hour,
		items:  make(map[string]data),
	}
}

func (s *Store) Create(adminID uint) string {
	id := randomID()
	s.mu.Lock()
	s.items[id] = data{AdminID: adminID, ExpiresAt: time.Now().Add(s.ttl)}
	s.mu.Unlock()
	return id
}

func (s *Store) GetAdminID(sessionID string) (uint, bool) {
	s.mu.RLock()
	item, ok := s.items[sessionID]
	s.mu.RUnlock()
	if !ok || time.Now().After(item.ExpiresAt) {
		if ok {
			s.Delete(sessionID)
		}
		return 0, false
	}
	return item.AdminID, true
}

func (s *Store) Delete(sessionID string) {
	s.mu.Lock()
	delete(s.items, sessionID)
	s.mu.Unlock()
}

func randomID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
