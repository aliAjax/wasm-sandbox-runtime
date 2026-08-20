package application

import (
	"context"
	"errors"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/wasm/inspector"
)

type AdmissionReview struct {
	ID              string    `json:"id"`
	ModuleVersionID string    `json:"module_version_id"`
	Risk            string    `json:"risk"`
	Approved        bool      `json:"approved"`
	Reasons         []string  `json:"reasons"`
	At              time.Time `json:"at"`
}

func ReviewModule(ctx context.Context, rt interface {
	Inspect(context.Context, []byte) (any, error)
}, v module.Version, b []byte) (AdmissionReview, error) {
	r := AdmissionReview{ID: NewID("review", v.ID), ModuleVersionID: v.ID, At: time.Now().UTC()}
	if len(b) > 64<<20 {
		return r, errors.New("module exceeds admission size")
	}
	raw, e := rt.Inspect(ctx, b)
	if e != nil {
		return r, e
	}
	info, ok := raw.(inspector.ModuleInfo)
	if !ok {
		return r, errors.New("inspector returned invalid type")
	}
	if len(info.Imports) > 8 {
		r.Risk = "high"
		r.Reasons = append(r.Reasons, "too many host imports")
	}
	if info.MemoryPages > 65536 {
		r.Risk = "high"
		r.Reasons = append(r.Reasons, "large linear memory")
	}
	if r.Risk == "" {
		r.Risk = "low"
		r.Approved = true
	}
	return r, nil
}
func (r AdmissionReview) Status() string {
	if r.Approved {
		return "approved"
	}
	return "quarantined"
}
