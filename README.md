# Block Indexer (Go, workspace)

Small Go services for indexing an EVM-like chain (plus a DAG overlay), exposing HTTP and WebSocket APIs. Everything is wired as Go modules in `go.work`, backed by Postgres and NATS, with Prometheus metrics on each service.

## Repository layout
- `services/` – individual services (see below) built with chi/zap/pgx.
- `pkg/` – shared helpers: config loader (Viper), Postgres pool + goose migrations, Zap logger, RPC stubs for EVM/DAG, Prometheus handler, NATS helper, and typed models.
- `configs/` – YAML config used by every service (`common.yaml` + one file per service) and `configs/prometheus.yml` for the bundled Prometheus container.
- `migrations/` – SQL migrations; `0001_init.sql` seeds core tables plus service-specific subfolders.
- `docker-compose.yml` – local stack: Postgres, NATS, Redis (currently unused by code), Prometheus, and all Go services launched with `go run`.
- `Makefile` – convenience targets to run each service with local Go.
- `go.work` / `go.work.sum` – Go 1.22 workspace tying the modules together.
- `docs/` – operational notes (cache ideas, data-service hints).

## Services
- `evm-indexer` (port 9001) – polls RPC stubs for new heads/backfill, stores blocks/tx/logs, exposes `/api/v1/internal/tip`, health, ready, metrics.
- `dag-indexer` (port 9002) – subscribes to DAG block events, tracks DAG tips, exposes `/api/v1/dag/tip`, health, ready, metrics.
- `api-gateway` (port 8080) – HTTP facade for blocks/tx/logs/dag/stats/search/contracts/traces; handlers rely on stub clients backed by Postgres connections.
- `stats-service` (port 9003) – periodic aggregation over `stats_daily`, endpoints for overview and historical stats.
- `trace-service` (port 9004) – queues and processes `debug_traceTransaction` calls, serves traces/state-diffs.
- `search-service` (port 9005) – subscribes to NATS events (placeholder) and serves `/api/v1/search`.
- `contract-service` (port 9006) – stores contract metadata and a stubbed `/verify` compiler hook.
- `ws-service` (port 8090) – WebSocket fanout at `/ws`; forwards NATS `evm.block.canonical` messages if present.

Most RPC calls and client responses are scaffolding/stubs; wire them to real chain nodes, event producers, and queries before production use.

## Configuration
- Common settings live in `configs/common.yaml` (service name, ports, Postgres, Redis placeholder, NATS subject prefix).
- Each service has its own config file in `configs/*.yaml` (e.g., `evm-indexer.yaml`, `dag-indexer.yaml`, `trace-service.yaml`).
- Config is loaded from `CONFIG_PATH` (default `./configs/common.yaml`) plus a service-specific env var:
  - `EVM_INDEXER_CONFIG`, `DAG_INDEXER_CONFIG`, `API_GATEWAY_CONFIG`, `STATS_SERVICE_CONFIG`, `TRACE_SERVICE_CONFIG`, `SEARCH_SERVICE_CONFIG`, `CONTRACT_SERVICE_CONFIG`, `WS_SERVICE_CONFIG`.
- Environment overrides use the service prefix (e.g., `EVM_INDEXER_HTTPPORT=9001`); see `.env.example` for all defaults.

## Running locally with Docker Compose
1) Copy and adjust envs: `cp .env.example .env` (Compose already points to `.env.example` if you want to edit in place).  
2) Start the stack:  
```bash
docker-compose up -d
```
   - Postgres on `localhost:5432`, NATS on `4222/8222`, Redis on `6379` (currently unused), Prometheus on `9090`.
   - App containers listen on the compose network only; add `ports:` if you need host access.

## Running services with the local Go toolchain
Requirements: Go 1.22+, Postgres + NATS reachable per `configs/common.yaml`.
- `make run-evm-indexer` (or any `run-*` target) starts a single service with `CONFIG_PATH` and the matching `*_CONFIG` set.
- `make build` builds all service binaries under `services/...`.
See `docs/evm-indexer-local.md` for a no-Docker loop when you want to iterate quickly on the EVM indexer only.

## Migrations
Migrations are not auto-run. Apply them with `goose` (or your preferred tool) against Postgres, e.g.:
```bash
POSTGRES_URL=postgres://explorer:explorer_pass@localhost:5432/explorer_db?sslmode=disable
goose -dir migrations/evm_indexer postgres "$POSTGRES_URL" up
goose -dir migrations/contract_service postgres "$POSTGRES_URL" up
```
Repeat for other folders (`dag_indexer`, `search_service`, `stats_service`, `trace_service`) as needed.

## Observability
- Each service exposes `/health`, `/ready`, and `/metrics` on its HTTP port.
- Prometheus in Compose uses `configs/prometheus.yml` to scrape every service.
- NATS is used for fanout plumbing; publishing from indexers is not wired yet, so search/ws will idle until events are produced.
