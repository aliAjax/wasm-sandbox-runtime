package worker

import (
	"crypto/sha256"
	"encoding/hex"
)

type CheckpointStore struct{ items map[string][]byte }

func NewCheckpointStore() *CheckpointStore { return &CheckpointStore{items: map[string][]byte{}} }
func (s *CheckpointStore) Put(id string, data []byte) string {
	s.items[id] = append([]byte(nil), data...)
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}
func (s *CheckpointStore) Get(id string) ([]byte, bool) {
	v, ok := s.items[id]
	return append([]byte(nil), v...), ok
}
