package application

import "sync"

type KeyLock struct {
	mu    sync.Mutex
	items map[string]*sync.Mutex
}

func NewKeyLock() *KeyLock { return &KeyLock{items: map[string]*sync.Mutex{}} }
func (l *KeyLock) Lock(key string) func() {
	l.mu.Lock()
	m := l.items[key]
	if m == nil {
		m = &sync.Mutex{}
		l.items[key] = m
	}
	l.mu.Unlock()
	m.Lock()
	return m.Unlock
}
