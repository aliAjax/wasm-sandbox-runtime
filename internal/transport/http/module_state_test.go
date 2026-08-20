package httptransport

import (
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"testing"
	"time"
)

func TestHTTPRecoveredModuleStatus(t *testing.T) {
	m := module.Module{ID: "m", State: module.Ready, Versions: []module.Version{{ID: "v", Digest: "d"}}}
	_, _ = m.Move(module.Quarantined, "scan", time.Now())
	_, _ = m.Move(module.Ready, "cleared", time.Now())
	if moduleStatus(m) != 200 {
		t.Fatal("HTTP status hid recovered module")
	}
}
func TestHTTPReadyModuleWithoutVersionIsUnavailable(t *testing.T) {
	if moduleStatus(module.Module{ID: "m", State: module.Ready}) != 409 {
		t.Fatal("module without a version was exposed as executable")
	}
}
func TestHTTPReadyModuleWithoutDigestIsUnavailable(t *testing.T) {
	if moduleStatus(module.Module{ID: "m", State: module.Ready, Versions: []module.Version{{ID: "v"}}}) != 409 {
		t.Fatal("module without a digest was exposed as executable")
	}
}
