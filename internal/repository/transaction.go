package repository

import (
	"context"
	"errors"
	"sync"
)

type UnitOfWork struct {
	mu      sync.Mutex
	closed  bool
	actions []func()
}

func Begin(_ context.Context) *UnitOfWork { return &UnitOfWork{} }
func (u *UnitOfWork) Add(action func()) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.closed {
		u.actions = append(u.actions, action)
	}
}
func (u *UnitOfWork) Commit() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return errors.New("transaction closed")
	}
	for _, a := range u.actions {
		a()
	}
	u.closed = true
	return nil
}
func (u *UnitOfWork) Rollback() { u.mu.Lock(); defer u.mu.Unlock(); u.closed = true; u.actions = nil }
