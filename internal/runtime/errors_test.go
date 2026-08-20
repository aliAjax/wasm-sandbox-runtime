package runtime

import (
	"errors"
	"testing"
)

func TestWrappedRuntimeErrorsPreserveChain(t *testing.T) {
	err := Wrap("timeout", "execute", ErrTimeout)
	if !errors.Is(err, ErrTimeout) {
		t.Fatal("wrapped timeout lost its error chain")
	}
	if !Retryable(err) || Classify(err) != "retryable" {
		t.Fatal("wrapped timeout was classified as terminal")
	}
}
func TestRuntimeErrorCodePreserved(t *testing.T) {
	if Wrap("timeout", "execute", ErrTimeout).(Error).CodeValue() != "timeout" {
		t.Fatal("runtime code was lost")
	}
}
