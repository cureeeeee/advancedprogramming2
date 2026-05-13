package store

import "sync"

type ProcessedStore struct {
	mu  sync.RWMutex
	ids map[string]struct{}
}

func NewProcessedStore() *ProcessedStore {
	return &ProcessedStore{ids: make(map[string]struct{})}
}

func (s *ProcessedStore) Seen(orderID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.ids[orderID]
	return ok
}

func (s *ProcessedStore) Mark(orderID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[orderID] = struct{}{}
}
