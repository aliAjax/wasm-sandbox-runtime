package httptransport

import (
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	tokens  float64
	updated time.Time
}
type RateLimiter struct {
	mu    sync.Mutex
	rate  float64
	burst float64
	items map[string]bucket
}

func NewRateLimiter(rate, burst float64) *RateLimiter {
	return &RateLimiter{rate: rate, burst: burst, items: map[string]bucket{}}
}
func (l *RateLimiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	b := l.items[key]
	if b.updated.IsZero() {
		b.updated = now
		b.tokens = l.burst
	}
	b.tokens += now.Sub(b.updated).Seconds() * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.updated = now
	if b.tokens < 1 {
		l.items[key] = b
		return false
	}
	b.tokens--
	l.items[key] = b
	return true
}
func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Tenant-ID")
		if key == "" {
			key = r.RemoteAddr
		}
		if !l.Allow(key, time.Now()) {
			w.Header().Set("Retry-After", "1")
			writeErr(w, rateError{})
			return
		}
		next.ServeHTTP(w, r)
	})
}

type rateError struct{}

func (rateError) Error() string { return "rate limit exceeded" }
