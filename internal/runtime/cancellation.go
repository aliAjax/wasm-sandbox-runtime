package runtime

import (
	"context"
	"sync"
)

type CancellationRegistry struct {
	mu    sync.Mutex
	items map[string]context.CancelFunc
}

func NewCancellationRegistry() *CancellationRegistry {
	return &CancellationRegistry{items: map[string]context.CancelFunc{}}
}
func (r *CancellationRegistry) Register(id string, parent context.Context) (context.Context, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; ok {
		return nil, context.Canceled
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.items[id] = cancel
	return ctx, nil
}
func (r *CancellationRegistry) Cancel(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	cancel, ok := r.items[id]
	if ok {
		cancel()
		delete(r.items, id)
	}
	return ok
}
func (r *CancellationRegistry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
}
