package auditevidence

import (
	runtimeport "github.com/example/wasm-sandbox-runtime/internal/runtime"
	"testing"
	"time"
)

func TestAuditVerifyUsesStableSnapshot(t *testing.T) {
	a := &runtimeport.AuditChain{}
	a.Append("e1", "run", "d1", time.Now())
	if !a.Verify() {
		t.Fatal("audit chain failed")
	}
}
