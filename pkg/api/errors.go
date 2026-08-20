package api

import "encoding/json"

type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Details   any    `json:"details,omitempty"`
}

func (e Error) Error() string       { return e.Code + ": " + e.Message }
func EncodeError(v any) []byte      { b, _ := json.Marshal(v); return b }
func Invalid(message string) Error  { return Error{Code: "invalid_argument", Message: message} }
func NotFound(message string) Error { return Error{Code: "not_found", Message: message} }
func Conflict(message string) Error { return Error{Code: "conflict", Message: message} }
func Internal(message string) Error { return Error{Code: "internal", Message: message} }
