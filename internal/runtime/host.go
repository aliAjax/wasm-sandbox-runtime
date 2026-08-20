package runtime

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Host struct {
	mu            sync.Mutex
	now           time.Time
	deterministic bool
	objects       map[string][]byte
	allowed       map[string]bool
}

func NewHost(now time.Time, deterministic bool) *Host {
	return &Host{now: now, deterministic: deterministic, objects: map[string][]byte{}, allowed: map[string]bool{}}
}
func (h *Host) Clock() time.Time        { h.mu.Lock(); defer h.mu.Unlock(); return h.now }
func (h *Host) Advance(d time.Duration) { h.mu.Lock(); defer h.mu.Unlock(); h.now = h.now.Add(d) }
func (h *Host) Random(out []byte) error {
	if h.deterministic {
		return errors.New("random disabled")
	}
	for i := range out {
		out[i] = byte((i*31 + 17) % 251)
	}
	return nil
}
func (h *Host) ReadObject(_ context.Context, key string) ([]byte, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	v, ok := h.objects[key]
	if !ok {
		return nil, errors.New("object not found")
	}
	return append([]byte(nil), v...), nil
}
func (h *Host) PutObject(key string, value []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.objects[key] = append([]byte(nil), value...)
}
func (h *Host) NetworkAllowed(host string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.allowed[host]
}
func (h *Host) AllowNetwork(host string) { h.mu.Lock(); defer h.mu.Unlock(); h.allowed[host] = true }
