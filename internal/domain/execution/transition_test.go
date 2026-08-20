package execution

import (
	"testing"
	"time"
)

func TestCheckpointRecoveryReachesSucceeded(t *testing.T) {
	e := Execution{State: Checkpointed, Attempts: []Attempt{{StartedAt: time.Now()}}}
	if err := e.FinishAttempt("sha256:ok", time.Now()); err != nil {
		t.Fatal(err)
	}
	if e.State != Succeeded || !e.HasResult() {
		t.Fatalf("state=%s result=%v", e.State, e.HasResult())
	}
}
func TestRecoveredSummaryIsTerminal(t *testing.T) {
	if !Summarize(Execution{State: Succeeded}).Terminal() {
		t.Fatal("succeeded summary was not terminal")
	}
}
func TestRecoveredFilterIncludesSucceeded(t *testing.T) {
	if !(Filter{States: []State{Succeeded}}).Recovered(Execution{State: Succeeded}) {
		t.Fatal("recovered state was filtered out")
	}
}
func TestSuccessfulAttemptClearsPreviousError(t *testing.T) {
	now := time.Now()
	e := Execution{State: Running, Error: "old", Attempts: []Attempt{{StartedAt: now}}}
	if err := e.FinishAttempt("sha256:ok", now); err != nil || e.Error != "" {
		t.Fatal("previous attempt error remained")
	}
}
