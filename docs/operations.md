# Operations

The in-memory adapter is intended for local validation. Production deployments replace repository ports with transactional SQL/object adapters and use a process-isolated WASM runtime. Keep the API behind TLS/mTLS or OIDC gateway, restrict egress through a proxy, and store secrets as Secret references. During graceful shutdown the scheduler stops admissions, waits for active workers, and HTTP server drains with a bounded timeout.

Run `scripts/smoke.sh` against a local process after startup. Stop the process with SIGINT once the smoke flow is complete.
