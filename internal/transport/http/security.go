package httptransport

import (
	"net/http"
	"strings"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}
func MaxHeader(next http.Handler, max int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var n int64
		for k, v := range r.Header {
			n += int64(len(k))
			for _, x := range v {
				n += int64(len(x))
			}
		}
		if n > max {
			writeErr(w, rateError{})
			return
		}
		next.ServeHTTP(w, r)
	})
}
