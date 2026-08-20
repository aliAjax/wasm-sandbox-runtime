package auditevidence

import (
	runtimeport "github.com/example/wasm-sandbox-runtime/internal/runtime"
	"sync"
	"testing"
	"time"
)

func TestAuditVerifyUsesStableSnapshot(t *testing.T) {
	a := &runtimeport.AuditChain{}
	a.Append("e1", "run", "d1", time.Now())
	start := make(chan struct{})
	var wg sync.WaitGroup
	ok := make(chan bool, 1)
	wg.Add(2)
	go func() { defer wg.Done(); <-start; a.Append("e2", "run", "d2", time.Now()) }()
	go func() { defer wg.Done(); <-start; ok <- a.Verify() }()
	close(start)
	wg.Wait()
	if !<-ok {
		t.Fatal("audit chain failed")
	}
}
func TestAuditRecordsCopied(t *testing.T) {
	a := &runtimeport.AuditChain{}
	a.Append("e1", "run", "d1", time.Now())
	start := make(chan struct{})
	out := make(chan []runtimeport.AuditRecord, 1)
	go func() { <-start; out <- a.StableRecords() }()
	go func() { <-start; a.Append("e2", "run", "d2", time.Now()) }()
	close(start)
	if len(<-out) == 0 {
		t.Fatal("empty audit snapshot")
	}
}
func TestMeterSnapshotStable(t *testing.T) {
	m := runtimeport.NewMeter(100, 100, 100)
	_ = m.ChargeInstructions(10)
	_ = m.ObserveMemory(20)
	_ = m.ChargeOutput(30)
	start := make(chan struct{})
	result := make(chan [3]uint64, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); <-start; _ = m.ChargeInstructions(1) }()
	go func() { defer wg.Done(); <-start; i, mem, out := m.SnapshotStable(); result <- [3]uint64{i, mem, out} }()
	close(start)
	wg.Wait()
	v := <-result
	if (v[0] != 10 && v[0] != 11) || v[1] != 20 || v[2] != 30 {
		t.Fatalf("misordered meter snapshot: %d %d %d", v[0], v[1], v[2])
	}
}
func TestTelemetryStrictBoundary(t *testing.T) {
	now := time.Now()
	td := runtimeport.Telemetry{ExecutionID: "e", Runtime: "r", StartedAt: now, FinishedAt: now}
	start := make(chan struct{})
	result := make(chan bool, 2)
	go func() { <-start; result <- td.ValidStrict() }()
	go func() { <-start; result <- td.ValidStrict() }()
	close(start)
	if !<-result || !<-result {
		t.Fatal("bad telemetry boundary")
	}
}
