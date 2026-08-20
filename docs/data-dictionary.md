# Data dictionary

| Field | Meaning |
|---|---|
| `tenant_id` | Tenant isolation key; required on module, execution and quota. |
| `module_version_id` | Immutable executable version identity. |
| `digest` | SHA-256 representation of module, input or output bytes. |
| `state` | Lifecycle state controlled by domain transition methods. |
| `idempotency_key` | Tenant-scoped duplicate submission key. |
| `lease_id` | Worker ownership token with expiry and optimistic version. |
| `limits` | CPU instruction, linear memory, output and timeout budgets. |
