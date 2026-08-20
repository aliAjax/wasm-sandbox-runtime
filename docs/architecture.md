# Architecture

Transport owns HTTP decoding, request limits, SSE and error mapping. Application coordinates use cases and emits state events. Domain packages contain lifecycle invariants, capability policy, module schemas and quota reservation without transport or storage imports. Repository ports isolate persistence. Worker owns fair queue, worker leases and recovery. Runtime ports isolate a safe simulator from a future Wasmtime/Wazero adapter.

```mermaid
sequenceDiagram
  participant C as Client
  participant A as API
  participant S as Service
  participant Q as Scheduler
  participant R as RuntimeAdapter
  C->>A: POST /executions
  A->>S: Submit with idempotency key
  S-->>C: 202 queued
  Q->>S: Admit + lease
  Q->>R: Execute with Meter and Policy
  R-->>S: digest/result or timeout
  S-->>A: state event
  A-->>C: SSE state
```
