package application

import (
	"errors"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/domain/policy"
)

func AuthorizeVersion(p policy.Policy, v module.Version) error {
	for _, perm := range v.EffectivePermissions() {
		switch perm {
		case module.PermNetwork:
			if !p.Network.AllowEgress {
				return errors.New("network capability not granted")
			}
		case module.PermFilesystem:
			if len(p.Mounts) == 0 {
				return errors.New("filesystem capability has no mount")
			}
		case module.PermRandom:
			if p.Deterministic {
				return errors.New("random capability conflicts with deterministic mode")
			}
		}
	}
	return nil
}
func EffectiveLimits(p policy.Policy, v module.Version) module.ResourceLimits {
	r := v.Limits
	if p.Limits.CPU > 0 && p.Limits.CPU < r.CPUInstructions {
		r.CPUInstructions = p.Limits.CPU
	}
	if p.Limits.Memory > 0 && p.Limits.Memory < r.MemoryBytes {
		r.MemoryBytes = p.Limits.Memory
	}
	if p.Limits.Output > 0 && p.Limits.Output < r.OutputBytes {
		r.OutputBytes = p.Limits.Output
	}
	return r
}
