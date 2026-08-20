package application

import (
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"testing"
	"time"
)

func TestServiceSeesRecoveredModule(t *testing.T) {
	m := module.Module{ID: "m", State: module.Ready}
	_, _ = m.Move(module.Quarantined, "scan", time.Now())
	_, _ = m.Move(module.Ready, "cleared", time.Now())
	if !ModuleExecutable(m) {
		t.Fatal("service rejected recovered module")
	}
}
