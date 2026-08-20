package runtime

import (
	"context"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/domain/policy"
)

type Request struct {
	Module        module.Version
	Input         []byte
	Policy        policy.Compiled
	Limits        module.ResourceLimits
	Deterministic bool
}
type Result struct {
	Output       []byte
	Digest       string
	Instructions uint64
	MemoryPeak   uint64
	Logs         []string
	StartedAt    time.Time
	FinishedAt   time.Time
}
type Checkpoint struct {
	Sequence uint64
	State    []byte
}
type RuntimeAdapter interface {
	Name() string
	Inspect(context.Context, []byte) (any, error)
	Execute(context.Context, Request) (Result, error)
	Cancel(string) error
	Checkpoint(context.Context, string) (Checkpoint, error)
	Restore(context.Context, string, Checkpoint) error
}
type HostContext interface {
	Clock() time.Time
	Random([]byte) error
	ReadObject(context.Context, string) ([]byte, error)
	NetworkAllowed(string) bool
}
