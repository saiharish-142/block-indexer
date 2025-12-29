GOFMT_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: fmt vet test build run-evm-indexer run-dag-indexer run-api-gateway run-stats-service run-trace-service run-search-service run-contract-service run-ws-service up down

fmt:
	gofmt -w $(GOFMT_FILES)

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./services/evm-indexer
	go build ./services/dag-indexer
	go build ./services/api-gateway
	go build ./services/stats-service
	go build ./services/trace-service
	go build ./services/search-service
	go build ./services/contract-service
	go build ./services/ws-service

run-evm-indexer:
	CONFIG_PATH=./configs/common.yaml EVM_INDEXER_CONFIG=./configs/evm-indexer.yaml go run ./services/evm-indexer

run-dag-indexer:
	CONFIG_PATH=./configs/common.yaml DAG_INDEXER_CONFIG=./configs/dag-indexer.yaml go run ./services/dag-indexer

run-api-gateway:
	CONFIG_PATH=./configs/common.yaml API_GATEWAY_CONFIG=./configs/api-gateway.yaml go run ./services/api-gateway

run-stats-service:
	CONFIG_PATH=./configs/common.yaml STATS_SERVICE_CONFIG=./configs/stats-service.yaml go run ./services/stats-service

run-trace-service:
	CONFIG_PATH=./configs/common.yaml TRACE_SERVICE_CONFIG=./configs/trace-service.yaml go run ./services/trace-service

run-search-service:
	CONFIG_PATH=./configs/common.yaml SEARCH_SERVICE_CONFIG=./configs/search-service.yaml go run ./services/search-service

run-contract-service:
	CONFIG_PATH=./configs/common.yaml CONTRACT_SERVICE_CONFIG=./configs/contract-service.yaml go run ./services/contract-service

run-ws-service:
	CONFIG_PATH=./configs/common.yaml WS_SERVICE_CONFIG=./configs/ws-service.yaml go run ./services/ws-service

up:
	docker-compose up -d

down:
	docker-compose down
