#!/usr/bin/env sh
set -eu
base="${BASE_URL:-http://localhost:8087}"
curl -fsS "$base/healthz" >/dev/null
curl -fsS "$base/readyz" >/dev/null
curl -fsS "$base/api/v1/runtime" >/dev/null
module=$(curl -fsS -X POST "$base/api/v1/modules" -H 'content-type: application/json' -d '{"tenant_id":"smoke","name":"smoke-module"}')
id=$(printf '%s' "$module" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$id"
curl -fsS "$base/api/v1/modules/$id" >/dev/null
echo "smoke ok"
