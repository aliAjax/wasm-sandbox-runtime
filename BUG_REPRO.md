# Bug Reproduction

Artifact read APIs expose shared slice storage. Mutating returned manifest chunks or reader snapshots changes later reads while stored digest strings remain unchanged.

```bash
go test ./pkg/artifact -run '^TestManifestReadIsolated$' -count=1
go test ./pkg/artifact -run '^TestCloneManifestIsolated$' -count=1
go test ./pkg/artifact -run '^TestReaderSnapshotIsolated$' -count=1
go test ./pkg/artifact -run '^TestChunksCopyIsolated$' -count=1
```

The buggy branch fails with messages including `manifest aliased`, `clone aliased`, `reader snapshot aliased`, and `chunk copy aliased`.
