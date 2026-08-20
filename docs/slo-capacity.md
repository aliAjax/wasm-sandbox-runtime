# SLO and capacity

默认目标为 99.9% API availability、执行提交 p99 < 500ms、状态 SSE 首事件 p99 < 1s。单实例默认 4 workers、256 队列槽位和 4 MiB 请求上限；按每任务 1 MiB 输入、10ms 模拟执行估算，每 worker 约 100 req/s。生产容量需依据 runtime 指令计量和 memory peak 压测校准。
