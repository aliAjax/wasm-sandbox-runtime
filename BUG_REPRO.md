# Bug Reproduction

Scheduler shutdown can return while workers remain active, accept work after closure, continue retries after cancellation, or panic when an empty priority queue is consumed.

```bash
go test -race ./internal/worker -run '^TestSchedulerStopDrainsWorkers$' -count=1
go test -race ./internal/worker -run '^TestRecoveryStopsOnContext$' -count=1
go test -race ./internal/worker -run '^TestPriorityQueueEmptyIsSafe$' -count=1
```

The buggy branch fails these tests with undrained workers, accepted work after stop, retries after cancellation, or an empty-heap panic.
