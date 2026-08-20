package runtime

import "context"

type AdapterFactory interface {
	Create(version string) (RuntimeAdapter, error)
}
type StaticFactory struct{ Adapter RuntimeAdapter }

func (f StaticFactory) Create(_ string) (RuntimeAdapter, error) { return f.Adapter, nil }
func ExecuteWithTimeout(ctx context.Context, adapter RuntimeAdapter, req Request) (Result, error) {
	return adapter.Execute(ctx, req)
}
