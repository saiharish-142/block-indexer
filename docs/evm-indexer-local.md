Run evm-indexer locally (no Docker)
===================================

This path runs the Go binary directly so you can iterate without waiting on Docker builds.

Prereqs
- Go 1.22+
- Postgres reachable on `localhost:5432`
- NATS reachable on `localhost:4222`
- `goose` CLI for migrations: `go install github.com/pressly/goose/v3/cmd/goose@latest`

1) Start Postgres and NATS
- Start Postgres however you prefer (Homebrew service, `pg_ctl`, etc.).
- Start NATS in another terminal: `nats-server -p 4222 -m 8222`.

2) Prepare the database
```bash
psql -h localhost -U postgres -c "CREATE USER explorer WITH PASSWORD 'explorer_pass';"
psql -h localhost -U postgres -c "CREATE DATABASE explorer_db OWNER explorer;"
```

3) Create a local config
- Copy `configs/common.yaml` to `configs/common.local.yaml` and change the hosts to your local services:
  - `db.host: localhost` (port `5432`)
  - `nats.url: nats://127.0.0.1:4222`
- Keep `configs/evm-indexer.yaml` as-is or point `rpc.evmEndpoint` / `rpc.wsEndpoint` to your test chain. You can also override any value with env vars such as `EVM_INDEXER_RPC_EVMENDPOINT=http://localhost:8545`.

4) Run migrations
```bash
POSTGRES_URL=postgres://explorer:explorer_pass@localhost:5432/explorer_db?sslmode=disable
goose -dir migrations/evm_indexer postgres "$POSTGRES_URL" up
```

5) Run the service (no Docker)
```bash
CONFIG_PATH=./configs/common.local.yaml \
EVM_INDEXER_CONFIG=./configs/evm-indexer.yaml \
go run ./services/evm-indexer
```
- If you do not want a new config file, set env overrides instead: `EVM_INDEXER_DB_HOST=localhost EVM_INDEXER_NATS_URL=nats://127.0.0.1:4222 go run ./services/evm-indexer`.

6) Quick checks
- `curl http://localhost:9001/health` (or `/ready` / `/api/v1/internal/tip`) should respond once the process is up.
- Stop with `Ctrl+C`; only the Go process restarts between code changes, keeping the loop fast.
