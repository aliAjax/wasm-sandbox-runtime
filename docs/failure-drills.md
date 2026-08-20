# Failure drills

- 终止 worker：确认 lease 过期、任务重新 admitted，幂等键仍返回原 execution。
- 注入超时：确认 `timed_out`、quota reservation release 和 SSE 事件。
- 提交恶意 import：确认 admission quarantine，未调用 Execute。
- 填满对象输出：确认输出 stream 拒绝，execution failed 且审计链完整。
- 回拨系统时间：生产时钟探针应告警并暂停 deterministic proofs。
