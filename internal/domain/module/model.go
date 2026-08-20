package module

import (
	"errors"
	"strings"
	"time"
)

type State string

const (
	Draft       State = "draft"
	Ready       State = "ready"
	Quarantined State = "quarantined"
	Revoked     State = "revoked"
)

type Permission string

const (
	PermClock      Permission = "clock"
	PermRandom     Permission = "random"
	PermFilesystem Permission = "filesystem.read"
	PermNetwork    Permission = "network.egress"
)

type ResourceLimits struct {
	CPUInstructions uint64        `json:"cpu_instructions"`
	MemoryBytes     uint64        `json:"memory_bytes"`
	OutputBytes     uint64        `json:"output_bytes"`
	Timeout         time.Duration `json:"timeout"`
}

func (r ResourceLimits) Validate() error {
	if r.CPUInstructions == 0 || r.CPUInstructions > 1e12 {
		return errors.New("cpu instruction budget out of range")
	}
	if r.MemoryBytes < 64*1024 || r.MemoryBytes > 8<<30 {
		return errors.New("memory limit out of range")
	}
	if r.OutputBytes == 0 || r.OutputBytes > 1<<30 {
		return errors.New("output limit out of range")
	}
	return nil
}

type Schema struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func (s Schema) Validate() error {
	if strings.TrimSpace(s.Input) == "" || strings.TrimSpace(s.Output) == "" {
		return errors.New("input and output schemas are required")
	}
	return nil
}

type Module struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	State       State     `json:"state"`
	Versions    []Version `json:"versions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Version struct {
	ID          string         `json:"id"`
	ModuleID    string         `json:"module_id"`
	Number      int            `json:"number"`
	Digest      string         `json:"digest"`
	Entrypoint  string         `json:"entrypoint"`
	Permissions []Permission   `json:"permissions"`
	Limits      ResourceLimits `json:"limits"`
	Schema      Schema         `json:"schema"`
	Signature   Signature      `json:"signature"`
	Imports     []string       `json:"imports"`
	Exports     []string       `json:"exports"`
	CreatedAt   time.Time      `json:"created_at"`
}
type Signature struct {
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"key_id"`
	Value     string `json:"value"`
	Verified  bool   `json:"verified"`
}

func (m *Module) AddVersion(v Version) error {
	if m.State == Revoked {
		return errors.New("module revoked")
	}
	if v.ModuleID != m.ID {
		return errors.New("module mismatch")
	}
	for _, old := range m.Versions {
		if old.Number == v.Number || old.Digest == v.Digest {
			return errors.New("duplicate module version")
		}
	}
	m.Versions = append(m.Versions, v)
	if m.State == Draft {
		m.State = Ready
	}
	return nil
}
func (v Version) Validate() error {
	if v.ID == "" || v.ModuleID == "" || v.Number < 1 || v.Digest == "" || v.Entrypoint == "" {
		return errors.New("invalid module version")
	}
	if err := v.Limits.Validate(); err != nil {
		return err
	}
	return v.Schema.Validate()
}
