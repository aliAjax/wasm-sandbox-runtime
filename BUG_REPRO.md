# Bug Reproduction

Request cancellation is detached from sandbox execution, output streaming, and host object reads. Work can continue after a timeout and subsequent requests can observe the wrong lifecycle.

```bash
go test ./ctxprobe -run '^TestCancellationPropagatesToAdapter$' -count=1
go test ./ctxprobe -run '^TestOutputStreamStopsAfterCancel$' -count=1
go test ./ctxprobe -run '^TestHostReadHonorsCancellation$' -count=1
```

The buggy branch fails because cancellation does not reach the adapter, output is written after cancellation, or host reads ignore the canceled context.
