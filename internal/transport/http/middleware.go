package httptransport

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/observability"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mw ...Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 12)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				jsonWrite(w, 500, map[string]string{"code": "internal", "message": "panic recovered"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(observability.WithTrace(r.Context())))
	})
}
func Timeout(d time.Duration) Middleware {
	return func(next http.Handler) http.Handler { return http.TimeoutHandler(next, d, "request timeout") }
}
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Idempotency-Key,X-Request-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
