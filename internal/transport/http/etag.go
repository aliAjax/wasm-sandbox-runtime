package httptransport

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

func ETag(v any) string {
	sum := sha256.Sum256([]byte(fmtValue(v)))
	return `"` + hex.EncodeToString(sum[:8]) + `"`
}
func fmtValue(v any) string  { return string(jsonBytes(v)) }
func jsonBytes(v any) []byte { b, _ := json.Marshal(v); return b }
func Conditional(w http.ResponseWriter, r *http.Request, tag string) bool {
	w.Header().Set("ETag", tag)
	return r.Header.Get("If-None-Match") == tag
}
