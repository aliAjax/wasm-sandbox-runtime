package runtime

import (
	"bytes"
	"context"
	"testing"
)

type contextProbe struct{ seen context.Context }

func (p *contextProbe) Name() string                                 { return "probe" }
func (p *contextProbe) Inspect(context.Context, []byte) (any, error) { return nil, nil }
func (p *contextProbe) Execute(ctx context.Context, _ Request) (Result, error) {
	p.seen = ctx
	return Result{}, nil
}
func (p *contextProbe) Cancel(string) error { return nil }
func (p *contextProbe) Checkpoint(context.Context, string) (Checkpoint, error) {
	return Checkpoint{}, nil
}
func (p *contextProbe) Restore(context.Context, string, Checkpoint) error { return nil }

func TestCancellationPropagatesToAdapter(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	r := NewCancellationRegistry()
	ctx, err := r.Register("e1", parent)
	if err != nil {
		t.Fatal(err)
	}
	p := &contextProbe{}
	_, _ = ExecuteWithTimeout(ctx, p, Request{})
	if p.seen != ctx {
		t.Fatal("adapter did not receive the registered context")
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
	s := NewOutputStream(&bytes.Buffer{}, 10)
	if _, err := s.WriteContext(ctx, []byte("x")); err == nil {
		t.Fatal("cancelled stream accepted output")
	}
}
