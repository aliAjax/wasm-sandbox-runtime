package api

import (
	"encoding/json"
	"errors"
	runtimeport "github.com/example/wasm-sandbox-runtime/internal/runtime"
)

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
func FromRuntimeError(err error) Error {
	if errors.Is(err, runtimeport.ErrTimeout) {
		return Error{Code: "timeout", Message: err.Error()}
	}
	return Internal(err.Error())
}
func RetryableStatus(err error) int { return 500 }
