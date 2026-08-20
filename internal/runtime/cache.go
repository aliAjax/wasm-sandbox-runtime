package runtime

import "sync"

type ResultCache struct {
	mu    sync.RWMutex
	items map[string]Result
}

func NewResultCache() *ResultCache { return &ResultCache{items: map[string]Result{}} }
func (c *ResultCache) Get(key string) (Result, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}
func (c *ResultCache) Put(key string, v Result) { c.mu.Lock(); defer c.mu.Unlock(); c.items[key] = v }
func (c *ResultCache) Delete(key string)        { c.mu.Lock(); defer c.mu.Unlock(); delete(c.items, key) }
func (c *ResultCache) Clear()                   { c.mu.Lock(); defer c.mu.Unlock(); c.items = map[string]Result{} }
