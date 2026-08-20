package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type traceKey struct{}

func WithTrace(ctx context.Context) context.Context {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return context.WithValue(ctx, traceKey{}, hex.EncodeToString(b))
}
func TraceID(ctx context.Context) string { v, _ := ctx.Value(traceKey{}).(string); return v }
func Span(ctx context.Context, name string) SpanRecord {
	return SpanRecord{TraceID: TraceID(ctx), Name: name, Attributes: map[string]string{}}
}

type SpanRecord struct {
	TraceID    string
	Name       string
	Attributes map[string]string
}

func (s *SpanRecord) Set(k, v string) { s.Attributes[k] = v }
