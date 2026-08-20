package module

import (
	"testing"
	"time"
)

func TestRecoveredModuleVisibility(t *testing.T) {
	m := Module{ID: "m", State: Ready}
	_, _ = m.Move(Quarantined, "scan", time.Now())
	_, _ = m.Move(Ready, "cleared", time.Now())
	if !m.CanExecute() {
		t.Fatal("recovered module is hidden")
	}
}
func TestRecoveredLifecycleReady(t *testing.T) {
	m := Module{ID: "m", State: Ready}
	_, _ = m.Move(Quarantined, "scan", time.Now())
	_, _ = m.Move(Ready, "cleared", time.Now())
	if !m.LifecycleReady() {
		t.Fatal("lifecycle did not return ready")
	}
}
func TestReadyTransitionIsIdempotent(t *testing.T) {
	if err := Transition(Ready, Ready); err != nil {
		t.Fatalf("ready event replay was rejected: %v", err)
	}
}
