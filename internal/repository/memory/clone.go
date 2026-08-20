package memory

import (
	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
)

func cloneModule(m module.Module) module.Module {
	m.Versions = append([]module.Version(nil), m.Versions...)
	return m
}
func cloneExecution(e execution.Execution) execution.Execution {
	e.Attempts = append([]execution.Attempt(nil), e.Attempts...)
	e.Checkpoints = append([]execution.Checkpoint(nil), e.Checkpoints...)
	return e
}
