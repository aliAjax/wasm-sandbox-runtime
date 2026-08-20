package runtime

import "sync"

type StateStore struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func NewStateStore() *StateStore { return &StateStore{items: map[string][]byte{}} }
func (s *StateStore) Save(id string, state []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = append([]byte(nil), state...)
}
func (s *StateStore) Load(id string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.items[id]
	return append([]byte(nil), v...), ok
}
func (s *StateStore) Remove(id string) { s.mu.Lock(); defer s.mu.Unlock(); delete(s.items, id) }
