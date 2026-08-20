package moduleevidence

import (
	module "github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"testing"
	"time"
)

func TestRecoveredModuleVisibility(t *testing.T) {
	m := module.Module{ID: "m", State: module.Ready}
	_, _ = m.Move(module.Quarantined, "scan", time.Now())
	_, _ = m.Move(module.Ready, "cleared", time.Now())
	if !m.CanExecute() {
		t.Fatal("recovered module is hidden")
	}
}
