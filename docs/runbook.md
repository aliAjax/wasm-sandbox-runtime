# Runbook

1. `GET /healthz` 检查进程，`GET /readyz` 检查是否接收流量。
2. 若队列延迟升高，检查 `WORKERS`、租户 quota 和 `/metrics`，不要直接增加无限并发。
3. 若 worker 崩溃，租约过期后由恢复循环重新入队；检查 execution attempt 和 checkpoint digest。
4. 若模块被 quarantine，使用 inspector report 查看 denied imports、risk 和 memory estimate。
5. 部署回滚只替换 runtime image；数据库 migration 是向前兼容，状态事件可重新消费。
