package api

import (
	runtimeport "github.com/example/wasm-sandbox-runtime/internal/runtime"
	"testing"
)

func TestRuntimeFailureMapsRetryable(t *testing.T) {
	err := runtimeport.Wrap("timeout", "execute", runtimeport.ErrTimeout)
	if got := FromRuntimeError(err); got.Code != "timeout" {
		t.Fatalf("code = %q", got.Code)
	}
}
func TestRuntimeCancelMapsRetryable(t *testing.T) {
	err := runtimeport.Wrap("cancelled", "execute", runtimeport.ErrCancelled)
	if got := FromRuntimeError(err); got.Code != "cancelled" {
		t.Fatalf("code = %q", got.Code)
	}
}
func TestRetryableStatusUsesRuntimeClassification(t *testing.T) {
	if got := RetryableStatus(runtimeport.ErrTimeout); got != 503 {
		t.Fatalf("timeout status = %d", got)
	}
}
