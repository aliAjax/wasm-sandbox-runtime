package application

import (
	"errors"
	"strings"

	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
)

func ValidateSubmission(e execution.Execution) error {
	if strings.TrimSpace(e.TenantID) == "" || strings.TrimSpace(e.ModuleVersionID) == "" {
		return errors.New("tenant and module version required")
	}
	if err := execution.ValidateIdempotencyKey(e.IdempotencyKey); err != nil {
		return err
	}
	return nil
}
func ValidateVersion(v module.Version) error {
	if err := module.ValidateName(v.Entrypoint); err != nil && v.Entrypoint != "_start" {
		return errors.New("entrypoint invalid")
	}
	if err := v.ValidateDigest(); err != nil {
		return err
	}
	return v.Validate()
}
