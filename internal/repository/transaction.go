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
	checked []func() error
}

func Begin(_ context.Context) *UnitOfWork { return &UnitOfWork{} }
func (u *UnitOfWork) Add(action func()) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.closed {
		u.actions = append(u.actions, action)
	}
}
func (u *UnitOfWork) AddChecked(action func() error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.closed {
		u.checked = append(u.checked, action)
	}
}
func (u *UnitOfWork) CommitChecked() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return errors.New("transaction closed")
	}
	u.closed = true
	for _, action := range u.checked {
		if err := action(); err != nil {
			u.actions = nil
			u.checked = nil
			return err
		}
	}
	u.actions = nil
	u.checked = nil
	return nil
}
func (u *UnitOfWork) Commit() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return errors.New("transaction closed")
	}
	u.closed = true
	for _, a := range u.actions {
		a()
	}
	u.actions = nil
	u.checked = nil
	return nil
}
func (u *UnitOfWork) Rollback() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.closed = true
	u.actions = nil
	u.checked = nil
}
func (u *UnitOfWork) Failed(err error) error {
	u.Rollback()
	return err
}
