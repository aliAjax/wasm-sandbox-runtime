package mock

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/runtime"
	"github.com/example/wasm-sandbox-runtime/internal/wasm/inspector"
	"github.com/example/wasm-sandbox-runtime/pkg/artifact"
)

type Runtime struct {
	mu          sync.Mutex
	cancelled   map[string]struct{}
	checkpoints map[string]runtime.Checkpoint
}

func NewRuntime() *Runtime {
	return &Runtime{cancelled: map[string]struct{}{}, checkpoints: map[string]runtime.Checkpoint{}}
}
func (r *Runtime) Name() string                                     { return "go-safe-simulator" }
func (r *Runtime) Inspect(_ context.Context, b []byte) (any, error) { return inspector.Inspect(b) }
func (r *Runtime) Execute(ctx context.Context, req runtime.Request) (runtime.Result, error) {
	start := time.Now().UTC()
	if err := req.Limits.Validate(); err != nil {
		return runtime.Result{}, err
	}
	if req.Module.Entrypoint == "" {
		return runtime.Result{}, errors.New("entrypoint is required")
	}
	if len(req.Input) > int(req.Limits.OutputBytes) {
		return runtime.Result{}, errors.New("input exceeds output budget")
	}
	select {
	case <-ctx.Done():
		return runtime.Result{}, ctx.Err()
	default:
	}
	r.mu.Lock()
	_, cancelled := r.cancelled[req.Module.ID]
	r.mu.Unlock()
	if cancelled {
		return runtime.Result{}, context.Canceled
	}
	out := transform(req.Input, req.Deterministic)
	if uint64(len(out)) > req.Limits.OutputBytes {
		return runtime.Result{}, errors.New("output limit exceeded")
	}
	finish := time.Now().UTC()
	return runtime.Result{Output: out, Digest: artifact.SHA256(out), Instructions: uint64(len(req.Input))*8 + 1, MemoryPeak: uint64(len(req.Input)), StartedAt: start, FinishedAt: finish, Logs: []string{"simulator execution completed"}}, nil
}
func transform(in []byte, deterministic bool) []byte {
	var v any
	if json.Unmarshal(in, &v) == nil {
		if deterministic {
			b, _ := json.Marshal(map[string]any{"ok": true, "value": v, "deterministic": true})
			return b
		}
		b, _ := json.Marshal(map[string]any{"ok": true, "value": v, "time": time.Now().UTC().UnixNano()})
		return b
	}
	out := append([]byte("wasm:"), in...)
	return out
}
func (r *Runtime) Cancel(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelled[id] = struct{}{}
	return nil
}
func (r *Runtime) Checkpoint(_ context.Context, id string) (runtime.Checkpoint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := runtime.Checkpoint{Sequence: uint64(len(r.checkpoints) + 1), State: []byte(id)}
	r.checkpoints[id] = cp
	return cp, nil
}
func (r *Runtime) Restore(_ context.Context, id string, cp runtime.Checkpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checkpoints[id] = cp
	return nil
}
