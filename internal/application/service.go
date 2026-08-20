package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/config"
	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/domain/policy"
	"github.com/example/wasm-sandbox-runtime/internal/domain/quota"
	"github.com/example/wasm-sandbox-runtime/internal/repository"
	runtimeport "github.com/example/wasm-sandbox-runtime/internal/runtime"
	"github.com/example/wasm-sandbox-runtime/pkg/clock"
)

type Service struct {
	store   repository.Store
	runtime runtimeport.RuntimeAdapter
	cfg     config.Config
	clock   clock.Clock
	log     *slog.Logger
	events  *EventHub
	mu      sync.Mutex
}

func NewService(s repository.Store, r runtimeport.RuntimeAdapter, c config.Config, l *slog.Logger) *Service {
	return &Service{store: s, runtime: r, cfg: c, clock: clock.System{}, log: l, events: NewEventHub()}
}
func (s *Service) Events() *EventHub { return s.events }
func (s *Service) CreateModule(ctx context.Context, m module.Module) (module.Module, error) {
	if m.ID == "" {
		m.ID = NewID("mod", m.Name)
	}
	if m.TenantID == "" || m.Name == "" {
		return m, errors.New("tenant_id and name required")
	}
	now := s.clock.Now()
	m.State = module.Draft
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.store.CreateModule(ctx, m); err != nil {
		return m, err
	}
	return m, nil
}
func (s *Service) AddVersion(ctx context.Context, moduleID string, v module.Version) (module.Version, error) {
	m, e := s.store.GetModule(ctx, moduleID)
	if e != nil {
		return v, e
	}
	if v.ID == "" {
		v.ID = NewID("ver", moduleID)
	}
	v.ModuleID = moduleID
	v.CreatedAt = s.clock.Now()
	if e = v.Validate(); e != nil {
		return v, e
	}
	if e = m.AddVersion(v); e != nil {
		return v, e
	}
	if e = s.store.PutVersion(ctx, v); e != nil {
		return v, e
	}
	return v, nil
}
func (s *Service) GetModule(ctx context.Context, id string) (module.Module, error) {
	return s.store.GetModule(ctx, id)
}
func (s *Service) ListModules(ctx context.Context, tenant string, offset, limit int) ([]module.Module, error) {
	return s.store.ListModules(ctx, tenant, offset, limit)
}
func (s *Service) Inspect(ctx context.Context, v module.Version, b []byte) (any, error) {
	info, e := s.runtime.Inspect(ctx, b)
	if e != nil {
		return nil, e
	}
	return info, nil
}
func (s *Service) PutPolicy(ctx context.Context, p policy.Policy) (policy.Policy, error) {
	if p.ID == "" {
		p.ID = NewID("pol", p.Name)
	}
	if p.Revision == 0 {
		p.Revision = 1
	}
	if e := p.Validate(); e != nil {
		return p, e
	}
	if e := s.store.PutPolicy(ctx, p); e != nil {
		return p, e
	}
	return p, nil
}
func (s *Service) GetPolicy(ctx context.Context, id string) (policy.Policy, error) {
	return s.store.GetPolicy(ctx, id)
}
func (s *Service) PutQuota(ctx context.Context, q quota.Quota) (quota.Quota, error) {
	if e := q.Validate(); e != nil {
		return q, e
	}
	if e := s.store.PutQuota(ctx, q); e != nil {
		return q, e
	}
	return q, nil
}
func (s *Service) GetQuota(ctx context.Context, id string) (quota.Quota, error) {
	return s.store.GetQuota(ctx, id)
}
func (s *Service) Submit(ctx context.Context, e execution.Execution) (execution.Execution, error) {
	if e.TenantID == "" || e.ModuleVersionID == "" {
		return e, errors.New("tenant_id and module_version_id required")
	}
	if e.IdempotencyKey != "" {
		if old, err := s.store.FindExecutionByKey(ctx, e.TenantID, e.IdempotencyKey); err == nil {
			return old, nil
		}
	}
	if e.ID == "" {
		e.ID = NewID("exe", e.ModuleVersionID)
	}
	e.State = execution.Queued
	e.CreatedAt = s.clock.Now()
	e.UpdatedAt = e.CreatedAt
	e.Deadline = e.CreatedAt.Add(s.cfg.DefaultTimeout)
	if e.Priority < 0 {
		e.Priority = 0
	}
	if e.Priority > 100 {
		e.Priority = 100
	}
	if err := s.store.CreateExecution(ctx, e); err != nil {
		return e, err
	}
	s.events.Publish(Event{ExecutionID: e.ID, State: string(e.State), At: e.CreatedAt})
	return e, nil
}
func (s *Service) GetExecution(ctx context.Context, id string) (execution.Execution, error) {
	return s.store.GetExecution(ctx, id)
}
func (s *Service) ListExecutions(ctx context.Context, tenant string, offset, limit int) ([]execution.Execution, error) {
	return s.store.ListExecutions(ctx, tenant, offset, limit)
}
func (s *Service) Admit(ctx context.Context, id string) error {
	e, err := s.store.GetExecution(ctx, id)
	if err != nil {
		return err
	}
	if err = e.Transition(execution.Admitted, "", s.clock.Now()); err != nil {
		return err
	}
	if err = s.store.UpdateExecution(ctx, e); err == nil {
		s.events.Publish(Event{ExecutionID: e.ID, State: string(e.State), At: e.UpdatedAt})
	}
	return err
}
func (s *Service) Run(ctx context.Context, id, workerID string) error {
	e, err := s.store.GetExecution(ctx, id)
	if err != nil {
		return err
	}
	v, err := s.store.GetVersion(ctx, e.ModuleVersionID)
	if err != nil {
		return err
	}
	p, err := s.store.GetPolicy(ctx, e.TenantID)
	if err != nil {
		p = policy.Policy{ID: "default", TenantID: e.TenantID, Limits: policy.Limits{CPU: v.Limits.CPUInstructions, Memory: v.Limits.MemoryBytes, Output: v.Limits.OutputBytes, TimeoutSeconds: int(v.Limits.Timeout / time.Second)}, Deterministic: true}
	}
	q, err := s.store.GetQuota(ctx, e.TenantID)
	if err != nil {
		return err
	}
	reservation, err := q.Reserve(float64(v.Limits.CPUInstructions)/1e6, v.Limits.MemoryBytes, v.Limits.OutputBytes, s.clock.Now())
	if err != nil {
		return err
	}
	defer func() { q.Release(reservation, s.clock.Now()); _ = s.store.UpdateQuota(context.Background(), q) }()
	if err = e.BeginAttempt(workerID, s.clock.Now()); err != nil {
		return err
	}
	if err = s.store.UpdateExecution(ctx, e); err != nil {
		return err
	}
	runCtx, cancel := context.WithDeadline(ctx, e.Deadline)
	defer cancel()
	res, runErr := s.runtime.Execute(runCtx, runtimeport.Request{Module: v, Input: []byte(`{"input":"sandbox"}`), Policy: p.Compile(), Limits: v.Limits, Deterministic: p.Deterministic})
	if runErr != nil {
		if errors.Is(runErr, context.DeadlineExceeded) {
			_ = e.Transition(execution.TimedOut, runErr.Error(), s.clock.Now())
		} else {
			_ = e.FailAttempt(runErr, s.clock.Now())
		}
	} else {
		_ = e.FinishAttempt(res.Digest, s.clock.Now())
	}
	if err = s.store.UpdateExecution(context.Background(), e); err == nil {
		s.events.Publish(Event{ExecutionID: e.ID, State: string(e.State), At: e.UpdatedAt, ResultDigest: e.ResultDigest})
	}
	return runErr
}
func (s *Service) Cancel(ctx context.Context, id string) error {
	e, err := s.store.GetExecution(ctx, id)
	if err != nil {
		return err
	}
	if execution.Terminal(e.State) {
		return errors.New("execution already terminal")
	}
	if err = s.runtime.Cancel(id); err != nil {
		return err
	}
	if err = e.Transition(execution.Cancelled, "cancelled by user", s.clock.Now()); err != nil {
		return err
	}
	if err = s.store.UpdateExecution(ctx, e); err == nil {
		s.events.Publish(Event{ExecutionID: e.ID, State: string(e.State), At: e.UpdatedAt})
	}
	return err
}
func (s *Service) Checkpoint(ctx context.Context, id string) (execution.Checkpoint, error) {
	e, err := s.store.GetExecution(ctx, id)
	if err != nil {
		return execution.Checkpoint{}, err
	}
	cp, err := s.runtime.Checkpoint(ctx, id)
	if err != nil {
		return execution.Checkpoint{}, err
	}
	sum := sha256.Sum256(cp.State)
	out := execution.Checkpoint{ID: NewID("cp", id), Sequence: cp.Sequence, Digest: "sha256:" + hex.EncodeToString(sum[:]), CreatedAt: s.clock.Now()}
	e.Checkpoints = append(e.Checkpoints, out)
	if e.State == execution.Running {
		_ = e.Transition(execution.Checkpointed, "checkpoint", out.CreatedAt)
	}
	err = s.store.UpdateExecution(ctx, e)
	return out, err
}
func NewID(prefix, seed string) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", prefix, seed, time.Now().UnixNano())))
	return prefix + "_" + hex.EncodeToString(h[:8])
}
