# Bug Reproduction

Concurrent audit export and append operations expose an unstable shared snapshot. Meter snapshot values are returned in the wrong positions and telemetry boundary checks are inconsistent.

```bash
go test -race ./auditchainprobe -run '^TestAuditVerifyUsesStableSnapshot$' -count=1
go test -race ./auditchainprobe -run '^TestAuditRecordsCopied$' -count=1
go test -race ./auditchainprobe -run '^TestMeterSnapshotStable$' -count=1
go test -race ./auditchainprobe -run '^TestTelemetryStrictBoundary$' -count=1
```

The buggy branch reports an audit-chain failure or data race in `internal/runtime/audit.go`, `misordered meter snapshot: 30 10 20`, and a bad telemetry boundary.
