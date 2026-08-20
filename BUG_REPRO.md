# Bug Reproduction

A module restored from quarantine can be `Ready` internally but remain unavailable to the application and HTTP layers. Replayed ready events and incomplete module metadata are also classified incorrectly.

```bash
go test ./internal/domain/module -run '^TestRecoveredModuleVisibility$' -count=1
go test ./internal/domain/module -run '^TestRecoveredLifecycleReady$' -count=1
go test ./internal/application -run '^TestServiceSeesRecoveredModule$' -count=1
go test ./internal/transport/http -run '^TestHTTPRecoveredModuleStatus$' -count=1
go test ./internal/domain/module -run '^TestReadyTransitionIsIdempotent$' -count=1
go test ./internal/transport/http -run '^TestHTTPReadyModuleWithoutVersionIsUnavailable$' -count=1
go test ./internal/transport/http -run '^TestHTTPReadyModuleWithoutDigestIsUnavailable$' -count=1
```

The buggy branch reports a hidden recovered module, a lifecycle that did not return ready, rejected execution, HTTP 409 for recovered state, rejected ready-event replay, or incomplete metadata exposed as executable.
