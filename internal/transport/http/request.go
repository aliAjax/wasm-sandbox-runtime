package httptransport

import (
	"encoding/json"
	"net/http"
)

type Envelope[T any] struct {
	Data      T      `json:"data"`
	RequestID string `json:"request_id,omitempty"`
}

func decodeStrict(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
func methodAllowed(r *http.Request, allowed ...string) bool {
	for _, m := range allowed {
		if r.Method == m {
			return true
		}
	}
	return false
}
func requireHeader(r *http.Request, name string) string { return r.Header.Get(name) }
