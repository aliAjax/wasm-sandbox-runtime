package application

import "sync"

type Idempotency struct {
	mu    sync.Mutex
	items map[string]string
}

func NewIdempotency() *Idempotency { return &Idempotency{items: map[string]string{}} }
func (i *Idempotency) Reserve(key, value string) (string, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if old, ok := i.items[key]; ok {
		return old, false
	}
	i.items[key] = value
	return value, true
}
func (i *Idempotency) Get(key string) (string, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v, ok := i.items[key]
	return v, ok
}
func (i *Idempotency) Delete(key string) { i.mu.Lock(); defer i.mu.Unlock(); delete(i.items, key) }
