package policy

import (
	"errors"
	"strings"
)

type Decision struct {
	Allowed    bool
	Reason     string
	Capability string
}

func (c Compiled) CheckCapability(capability string) Decision {
	if capability == "" {
		return Decision{Allowed: false, Reason: "empty capability"}
	}
	if capability == "clock" && c.Deterministic {
		return Decision{Allowed: true, Reason: "deterministic clock provided", Capability: capability}
	}
	if capability == "random" && c.Deterministic {
		return Decision{Allowed: false, Reason: "randomness disabled in deterministic mode", Capability: capability}
	}
	for _, d := range c.Domains {
		if strings.HasPrefix(capability, "network:") && strings.TrimPrefix(capability, "network:") == d {
			return Decision{Allowed: true, Reason: "domain allowlist match", Capability: capability}
		}
	}
	return Decision{Allowed: false, Reason: "capability not in compiled policy", Capability: capability}
}
func (p Policy) ValidateEnvironment(env map[string]string) error {
	for k := range env {
		if !p.AllowsEnv(k) {
			return errors.New("environment variable denied: " + k)
		}
	}
	return nil
}
func (p Policy) ValidateMount(name string) error {
	for _, m := range p.Mounts {
		if m.Name == name && m.ReadOnly {
			return nil
		}
	}
	return errors.New("mount denied: " + name)
}
