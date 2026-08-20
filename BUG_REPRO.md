# Bug Reproduction

Failed transactional actions do not close and release the unit of work consistently. Checked commits can swallow an error, later actions can still run, and completed actions can remain retained.

```bash
go test ./internal/repository -run '^TestUnitOfWorkRollsBackFailedAction$' -count=1
go test ./internal/repository -run '^TestRollbackReleasesActions$' -count=1
go test ./internal/repository -run '^TestCommitCheckedStopsOnError$' -count=1
go test ./internal/repository -run '^TestStoreContractRejectsNil$' -count=1
go test ./internal/repository -run '^TestCommitReleasesCompletedActions$' -count=1
```

The buggy branch reports repeated or retained actions, continued execution after the first checked-action error, or acceptance of a nil store.
