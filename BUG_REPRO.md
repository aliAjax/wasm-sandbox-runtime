# Bug Reproduction

Concurrent event publication and subscription cancellation can race with channel closure, causing `send on closed channel`. Event snapshots can also expose or mix mutable attribute maps.

```bash
go test -race ./internal/application -run '^TestEventHubConcurrentUnsubscribe$' -count=1
go test -race ./internal/application -run '^TestEventSnapshotsAreIndependent$' -count=1
```

The buggy branch reports a race or panic in `internal/application/events.go`, and the snapshot isolation test reports mutated or mismatched attributes.
