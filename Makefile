PROJECT=github.com/example/block-indexer

.PHONY: build test lint docker-build run-local migrate

build:
	GOOS=linux GOARCH=amd64 go build ./cmd/...

test:
	go test ./... -race -cover

lint:
	golangci-lint run ./...

docker-build:
	docker build -f deploy/docker/Dockerfile.indexer -t block-indexer-indexer:local .
	docker build -f deploy/docker/Dockerfile.api -t block-indexer-api:local .
	docker build -f deploy/docker/Dockerfile.ws -t block-indexer-ws:local .

run-local:
	APP_ENV=development go run ./cmd/api

migrate:
	goose -dir ./migrations postgres "$$POSTGRES_URL" up
