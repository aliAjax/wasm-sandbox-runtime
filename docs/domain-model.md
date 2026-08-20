# Domain model

`Module` 聚合拥有多个不可重复的 `Version`。版本包含入口、导入导出、权限、签名、JSON schema 和资源限制。`Execution` 记录输入摘要、状态机、租约、attempt 和 checkpoint。`Policy` 编译成不可变 capability set；`Quota` 通过 Reservation 预占并以幂等 Release 回滚。

输入只接受摘要或受限字节流。对象端口使用分块 manifest 和 SHA-256，服务不会管理用户文件本体。
