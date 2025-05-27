package mem

import (
	"errors"
	"sync"
	"time"

	"github.com/inmore/gopaste/internal/model"
	"github.com/inmore/gopaste/internal/storage"
)

var _ storage.Storage = (*Store)(nil)

type Store struct {
	mu   sync.RWMutex
	data map[string]*model.Paste
}

func New() *Store {
	return &Store{data: make(map[string]*model.Paste)}
}

func (s *Store) Save(p *model.Paste) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[p.ID] = p
	return nil
}

func (s *Store) Load(id string) (*model.Paste, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[id]
	if !ok || time.Now().After(p.ExpiresAt) {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (s *Store) DeleteExpired() (int, error) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	expiredKeys := make([]string, 0, len(s.data))
	for id, p := range s.data {
		if now.After(p.ExpiresAt) {
			expiredKeys = append(expiredKeys, id)
		}
	}
	for _, id := range expiredKeys {
		delete(s.data, id)
	}
	return len(expiredKeys), nil
}

func (s *Store) Close() error { return nil }
