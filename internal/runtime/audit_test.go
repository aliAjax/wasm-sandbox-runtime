package runtime

import (
	"sync"
	"testing"
	"time"
)

func TestAuditVerifyUsesStableSnapshot(t *testing.T) {
	a := &AuditChain{}
	a.Append("e1", "run", "d1", time.Now())
	if !a.Verify() {
		t.Fatal("audit chain failed")
	}
}
func TestAuditRecordsSnapshotIsolated(t *testing.T) {
	a := &AuditChain{}
	a.Append("e1", "run", "d1", time.Now())
	start := make(chan struct{})
	out := make(chan []AuditRecord, 1)
	go func() { <-start; out <- a.StableRecords() }()
	go func() { <-start; a.Append("e2", "run", "d2", time.Now()) }()
	close(start)
	v := <-out
	if len(v) == 0 {
		t.Fatal("empty audit snapshot")
	}
}
func TestMeterSnapshotStable(t *testing.T) {
	m := NewMeter(100, 100, 100)
	_ = m.ChargeInstructions(10)
	_ = m.ObserveMemory(20)
	_ = m.ChargeOutput(30)
	start := make(chan struct{})
	var wg sync.WaitGroup
	result := make(chan [3]uint64, 1)
	wg.Add(2)
	go func() { defer wg.Done(); <-start; _ = m.ChargeInstructions(1) }()
	go func() {
		defer wg.Done()
		<-start
		i, memory, output := m.SnapshotStable()
		result <- [3]uint64{i, memory, output}
	}()
	close(start)
	wg.Wait()
	v := <-result
	if (v[0] != 10 && v[0] != 11) || v[1] != 20 || v[2] != 30 {
		t.Fatalf("misordered meter snapshot: %d %d %d", v[0], v[1], v[2])
	}
}
func TestTelemetryStrictBoundary(t *testing.T) {
	now := time.Now()
	tdata := Telemetry{ExecutionID: "e", Runtime: "r", StartedAt: now, FinishedAt: now}
	start := make(chan struct{})
	results := make(chan bool, 2)
	go func() { <-start; results <- tdata.ValidStrict() }()
	go func() { <-start; results <- tdata.ValidStrict() }()
	close(start)
	first, second := <-results, <-results
	if !first || !second {
		t.Fatal("bad telemetry boundary")
	}
}
