# Threat model

攻击面包括畸形 WASM section、恶意 imports、超大线性内存、输出洪泛、租户配额耗尽、SSE 信息泄露和重复提交。检查器先验证 magic、版本、section 长度和 LEB128 边界，不实例化模块；策略编译器拒绝未声明能力、可写挂载和 deterministic 随机源；Meter、输出流和 MaxBytesReader 提供三重预算。审计链记录状态和摘要，不写入敏感输入。

运行在非 root、只读根文件系统、无新增 Linux capability 的容器。生产适配器应将每个任务放入独立 runtime worker，并配置 seccomp、cgroup 和 egress proxy。
