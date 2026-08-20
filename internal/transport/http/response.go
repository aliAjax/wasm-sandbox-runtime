package httptransport

import (
	"encoding/json"
	"net/http"
)

type Problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func problem(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	json.NewEncoder(w).Encode(Problem{Type: "about:blank", Title: http.StatusText(status), Status: status, Detail: detail})
}
func accepted(w http.ResponseWriter, v any) { jsonWrite(w, http.StatusAccepted, v) }
func created(w http.ResponseWriter, v any)  { jsonWrite(w, http.StatusCreated, v) }
func noContent(w http.ResponseWriter)       { w.WriteHeader(http.StatusNoContent) }
