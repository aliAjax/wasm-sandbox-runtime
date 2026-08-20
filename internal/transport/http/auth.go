package httptransport

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Principal struct {
	TenantID string
	Subject  string
	Method   string
}

func Authenticate(r *http.Request) (Principal, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return Principal{TenantID: r.Header.Get("X-Tenant-ID"), Subject: "anonymous", Method: "none"}, nil
	}
	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return Principal{}, errors.New("invalid authorization header")
	}
	if len(parts[1]) < 8 {
		return Principal{}, errors.New("token too short")
	}
	return Principal{TenantID: r.Header.Get("X-Tenant-ID"), Subject: parts[1][:8], Method: "bearer"}, nil
}
func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, e := Authenticate(r)
		if e != nil {
			writeErr(w, e)
			return
		}
		if p.TenantID == "" {
			writeErr(w, errors.New("tenant header required"))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, p)))
	})
}

type principalKey struct{}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}
