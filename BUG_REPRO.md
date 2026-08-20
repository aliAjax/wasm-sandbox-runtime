# Bug Reproduction

A checkpointed execution cannot complete cleanly: recovery can remain outside `Succeeded`, preserve an earlier error, and disagree with summary and recovered-list views.

```bash
go test ./stateprobe -run '^TestCheckpointRecoveryReachesSucceeded$' -count=1
go test ./stateprobe -run '^TestRecoveredSummaryIsTerminal$' -count=1
go test ./stateprobe -run '^TestRecoveredFilterIncludesSucceeded$' -count=1
go test ./stateprobe -run '^TestSuccessfulAttemptClearsPreviousError$' -count=1
```

On the buggy branch all four probes fail, reporting an invalid checkpoint transition, a non-terminal recovered summary, a missing recovered execution, or a stale previous error.
