package memory

import (
	"context"
	"sync"

	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/domain/policy"
	"github.com/example/wasm-sandbox-runtime/internal/domain/quota"
	"github.com/example/wasm-sandbox-runtime/internal/repository"
)

type Store struct {
	mu         sync.RWMutex
	modules    map[string]module.Module
	versions   map[string]module.Version
	policies   map[string]policy.Policy
	quotas     map[string]quota.Quota
	executions map[string]execution.Execution
	keys       map[string]string
}

func NewStore() *Store {
	return &Store{modules: map[string]module.Module{}, versions: map[string]module.Version{}, policies: map[string]policy.Policy{}, quotas: map[string]quota.Quota{}, executions: map[string]execution.Execution{}, keys: map[string]string{}}
}
func (s *Store) CreateModule(_ context.Context, m module.Module) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.modules[m.ID]; ok {
		return repository.ErrNotFound
	}
	s.modules[m.ID] = m
	return nil
}
func (s *Store) GetModule(_ context.Context, id string) (module.Module, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.modules[id]
	if !ok {
		return module.Module{}, repository.ErrNotFound
	}
	m.Versions = append([]module.Version(nil), m.Versions...)
	return m, nil
}
func (s *Store) ListModules(_ context.Context, tenant string, offset, limit int) ([]module.Module, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]module.Module, 0)
	for _, m := range s.modules {
		if tenant == "" || m.TenantID == tenant {
			all = append(all, m)
		}
	}
	if offset > len(all) {
		offset = len(all)
	}
	end := offset + limit
	if limit <= 0 || end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}
func (s *Store) PutVersion(_ context.Context, v module.Version) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.versions[v.ID]; ok {
		return repository.ErrNotFound
	}
	s.versions[v.ID] = v
	m, ok := s.modules[v.ModuleID]
	if !ok {
		return repository.ErrNotFound
	}
	m.Versions = append(m.Versions, v)
	s.modules[m.ID] = m
	return nil
}
func (s *Store) GetVersion(_ context.Context, id string) (module.Version, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.versions[id]
	if !ok {
		return module.Version{}, repository.ErrNotFound
	}
	return v, nil
}
func (s *Store) PutPolicy(_ context.Context, p policy.Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[p.ID] = p
	return nil
}
func (s *Store) GetPolicy(_ context.Context, id string) (policy.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.policies[id]
	if !ok {
		return policy.Policy{}, repository.ErrNotFound
	}
	return p, nil
}
func (s *Store) PutQuota(_ context.Context, q quota.Quota) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.quotas[q.TenantID] = q
	return nil
}
func (s *Store) GetQuota(_ context.Context, id string) (quota.Quota, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.quotas[id]
	if !ok {
		return quota.Quota{}, repository.ErrNotFound
	}
	return q, nil
}
func (s *Store) UpdateQuota(_ context.Context, q quota.Quota) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.quotas[q.TenantID] = q
	return nil
}
func (s *Store) CreateExecution(_ context.Context, e execution.Execution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.IdempotencyKey != "" {
		k := e.TenantID + ":" + e.IdempotencyKey
		if _, ok := s.keys[k]; ok {
			return repository.ErrNotFound
		}
		s.keys[k] = e.ID
	}
	s.executions[e.ID] = e
	return nil
}
func (s *Store) GetExecution(_ context.Context, id string) (execution.Execution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.executions[id]
	if !ok {
		return execution.Execution{}, repository.ErrNotFound
	}
	return e, nil
}
func (s *Store) UpdateExecution(_ context.Context, e execution.Execution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executions[e.ID]; !ok {
		return repository.ErrNotFound
	}
	s.executions[e.ID] = e
	return nil
}
func (s *Store) FindExecutionByKey(_ context.Context, tenant, key string) (execution.Execution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.keys[tenant+":"+key]
	if !ok {
		return execution.Execution{}, repository.ErrNotFound
	}
	return s.executions[id], nil
}
func (s *Store) ListExecutions(_ context.Context, tenant string, offset, limit int) ([]execution.Execution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := []execution.Execution{}
	for _, e := range s.executions {
		if tenant == "" || e.TenantID == tenant {
			all = append(all, e)
		}
	}
	if offset > len(all) {
		offset = len(all)
	}
	end := offset + limit
	if limit <= 0 || end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}
