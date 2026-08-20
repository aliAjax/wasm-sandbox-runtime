package contextevidence

import (
	"bytes"
	"context"
	"errors"
	runtimeport "github.com/example/wasm-sandbox-runtime/internal/runtime"
	"testing"
	"time"
)

type probe struct{ seen context.Context }

func (p *probe) Name() string                                 { return "probe" }
func (p *probe) Inspect(context.Context, []byte) (any, error) { return nil, nil }
func (p *probe) Execute(c context.Context, _ runtimeport.Request) (runtimeport.Result, error) {
	p.seen = c
	return runtimeport.Result{}, nil
}
func (p *probe) Cancel(string) error { return nil }
func (p *probe) Checkpoint(context.Context, string) (runtimeport.Checkpoint, error) {
	return runtimeport.Checkpoint{}, nil
}
func (p *probe) Restore(context.Context, string, runtimeport.Checkpoint) error { return nil }
func TestCancellationPropagatesToAdapter(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	r := runtimeport.NewCancellationRegistry()
	ctx, _ := r.Register("e1", parent)
	p := &probe{}
	_, _ = runtimeport.ExecuteWithTimeout(ctx, p, runtimeport.Request{})
	if p.seen != ctx {
		t.Fatal("adapter did not receive context")
	}
	cancel()
	select {
	case <-ctx.Done():
	default:
		t.Fatal("parent cancellation was lost")
	}
}
func TestOutputStreamStopsAfterCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := runtimeport.NewOutputStream(&bytes.Buffer{}, 10)
	if _, err := s.WriteContext(ctx, []byte("x")); err == nil {
		t.Fatal("cancelled stream accepted output")
	}
}
func TestHostReadHonorsCancellation(t *testing.T) {
	h := runtimeport.NewHost(time.Now(), true)
	h.PutObject("state", []byte("value"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := h.ReadObject(ctx, "state"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled object read returned %v", err)
	}
}
