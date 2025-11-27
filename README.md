# Block Explorer Backend (Go)

Production-oriented scaffold for a high-speed EVM-like chain (≈10 blocks/sec) built as Go microservices: indexer, public API, and WebSocket fanout. Includes gRPC contracts, Postgres/Redis data layer, Docker/K8s/Helm deployment, and observability hooks.

## Layout
- `cmd/indexer`: chain ingestion service (WS heads + polling backfill, bulk inserts).
- `cmd/api`: REST API (chi) with pagination stubs.
- `cmd/ws`: WebSocket service for heads/tx/address topics.
- `internal/*`: shared config, logging, metrics, telemetry, db/cache helpers, gRPC server glue.
- `protos/explorer.proto`: gRPC definitions (replace `internal/pb` placeholder with generated code).
- `migrations/`: Postgres schema with partitioned tables.
- `deploy/docker/`: Multi-stage Dockerfiles per service.
- `deploy/k8s/`: Minimal manifests for Deployments/Services/ConfigMap/Secret.
- `deploy/helm/`: Helm chart skeleton.
- `deploy/observability/`: Grafana dashboard stub.
- `docs/`: Cache strategy and data layer notes.

## Requirements
- Go 1.21+
- Docker
- Make
- (Optional) `golangci-lint`, `goose` for migrations

## Setup
```bash
go mod tidy          # download deps
gofmt -w .           # format
```

## Makefile targets
- `make build` – build all binaries under `cmd/...` (linux/amd64).
- `make test` – `go test ./... -race -cover`.
- `make lint` – runs `golangci-lint`.
- `make docker-build` – builds local images for indexer/api/ws using Dockerfiles.
- `make run-local` – runs API in dev mode (env-driven config).
- `make migrate` – applies migrations with goose (`POSTGRES_URL` must be set).

## Docker images
Example builds (without Make):
```bash
docker build -f deploy/docker/Dockerfile.indexer -t block-indexer-indexer:local .
docker build -f deploy/docker/Dockerfile.api     -t block-indexer-api:local .
docker build -f deploy/docker/Dockerfile.ws      -t block-indexer-ws:local .
```

## Kubernetes
- Apply config/secret (edit placeholders), then Deployments/Services:
```bash
kubectl apply -f deploy/k8s/configmap.yaml
kubectl apply -f deploy/k8s/secret-example.yaml   # replace with real secret management
kubectl apply -f deploy/k8s/indexer.yaml
kubectl apply -f deploy/k8s/api.yaml
kubectl apply -f deploy/k8s/ws.yaml
```
- Helm skeleton is under `deploy/helm/` (fill values, add secrets).

## Observability
- `/metrics` on each service via Prometheus client; key metrics defined in `internal/metrics`.
- OTEL tracer stub in `internal/telemetry` (wire collector endpoint in config/env).
- Grafana dashboard stub in `deploy/observability/grafana-dashboard.json`.

## Migrations
Using `goose` (or swap for `golang-migrate`). Example:
```bash
POSTGRES_URL=postgres://user:pass@localhost:5432/block_indexer?sslmode=disable \
goose -dir ./migrations postgres "$POSTGRES_URL" up
```

## Next steps
- Generate real gRPC code from `protos/explorer.proto` (buf or protoc) and replace `internal/pb`.
- Implement chain RPC logic (go-ethereum) and full DB/cache wiring with reorg handling.
- Harden configs (timeouts, retry/backoff), add integration tests, and tune partitions/indices.
