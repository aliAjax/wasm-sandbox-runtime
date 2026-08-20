# Bug Reproduction

Minimal configuration can write into nil maps and panic with `assignment to entry in nil map`. Missing or empty secret references may also be treated as valid values.

```bash
go test ./internal/config -run '^TestMinimalConfigIsSafe$' -count=1
go test ./internal/config -run '^TestMissingSecretIsRejected$' -count=1
```

The buggy branch panics during minimal configuration or fails because a missing secret is accepted instead of being rejected early.
