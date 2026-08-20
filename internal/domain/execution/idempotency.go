package execution

import (
	"errors"
	"strings"
)

func ValidateIdempotencyKey(k string) error {
	if strings.TrimSpace(k) == "" {
		return nil
	}
	if len(k) > 256 {
		return errors.New("idempotency key too long")
	}
	if strings.ContainsAny(k, "\r\n") {
		return errors.New("idempotency key contains control characters")
	}
	return nil
}
func SameRequest(a, b Execution) bool {
	return a.TenantID == b.TenantID && a.ModuleVersionID == b.ModuleVersionID && a.InputDigest == b.InputDigest
}
