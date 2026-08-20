package artifact

import (
	"errors"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	objects map[string]Manifest
	data    map[string][]byte
}

func NewStore() *Store { return &Store{objects: map[string]Manifest{}, data: map[string][]byte{}} }
func (s *Store) Put(data []byte, chunkSize int) Manifest {
	m := BuildManifest(data, chunkSize)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.objects[m.Digest] = m
	s.data[m.Digest] = append([]byte(nil), data...)
	return m
}
func (s *Store) Get(digest string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[digest]
	if !ok {
		return nil, errors.New("artifact not found")
	}
	return append([]byte(nil), v...), nil
}
func (s *Store) Has(digest string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.objects[digest]
	return ok
}
func (s *Store) Delete(digest string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.objects, digest)
	delete(s.data, digest)
}
func (s *Store) Manifest(digest string) (Manifest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.objects[digest]
	if !ok {
		return Manifest{}, errors.New("artifact not found")
	}
	return m, nil
}
func (m Manifest) ChunksCopy() []Chunk { return m.Chunks }
