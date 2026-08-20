package repository

import (
	"context"
	"errors"

	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/domain/policy"
	"github.com/example/wasm-sandbox-runtime/internal/domain/quota"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	CreateModule(context.Context, module.Module) error
	GetModule(context.Context, string) (module.Module, error)
	ListModules(context.Context, string, int, int) ([]module.Module, error)
	PutVersion(context.Context, module.Version) error
	GetVersion(context.Context, string) (module.Version, error)
	PutPolicy(context.Context, policy.Policy) error
	GetPolicy(context.Context, string) (policy.Policy, error)
	PutQuota(context.Context, quota.Quota) error
	GetQuota(context.Context, string) (quota.Quota, error)
	UpdateQuota(context.Context, quota.Quota) error
	CreateExecution(context.Context, execution.Execution) error
	GetExecution(context.Context, string) (execution.Execution, error)
	UpdateExecution(context.Context, execution.Execution) error
	FindExecutionByKey(context.Context, string, string) (execution.Execution, error)
	ListExecutions(context.Context, string, int, int) ([]execution.Execution, error)
}
