# Bug Reproduction

Wrapped runtime timeout and cancellation errors lose their identity and code, so retry classification fails and the API returns a generic 500 response.

```bash
go test ./internal/runtime -run '^TestWrappedRuntimeErrorsPreserveChain$' -count=1
go test ./internal/runtime -run '^TestRuntimeErrorCodePreserved$' -count=1
go test ./pkg/api -run '^TestRuntimeFailureMapsRetryable$' -count=1
go test ./pkg/api -run '^TestRuntimeCancelMapsRetryable$' -count=1
go test ./pkg/api -run '^TestRetryableStatusUsesRuntimeClassification$' -count=1
```

The buggy branch emits failures showing that `runtime_timeout: execute: runtime timeout` cannot be matched through the error chain and maps to status 500 instead of a retryable response.
