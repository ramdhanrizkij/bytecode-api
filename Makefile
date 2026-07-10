-include .env
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export

# ─── App ─────────────────────────────────────────────────────────────────────

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go
	go build -o bin/migrate cmd/migrate/main.go

run:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

# ─── Test ─────────────────────────────────────────────────────────────────────

test:
	go test -v ./...

test-unit:
	go test -v $$(go list ./... | grep -v /test/integration)

test-integration:
	go test -v -tags=integration ./test/integration/...

tidy:
	go mod tidy

# ─── API Documentation ──────────────────────────────────────────────────────
# Requires the `swag` CLI: https://github.com/swaggo/swag

swagger-setup:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	swag init -g cmd/api/main.go --parseDependency --parseInternal --output docs

# ─── Migrations ──────────────────────────────────────────────────────────────
# DATABASE_URL can be overridden on the command line:
#   make migrate-up DATABASE_URL=postgres://...

DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Apply all pending migrations
migrate-up:
	go run ./cmd/migrate -action up -database "$(DATABASE_URL)"

# Roll back the last applied migration
migrate-down:
	go run ./cmd/migrate -action down -database "$(DATABASE_URL)"

# Drop all migrations and apply them from scratch
migrate-refresh:
	go run ./cmd/migrate -action refresh -database "$(DATABASE_URL)"

# Create a new Go migration
# Usage: make migrate-create name=create_something_table
migrate-create:
	@test -n "$(name)" || (echo "ERROR: name is required. Usage: make migrate-create name=create_xxx_table" && exit 1)
	go run ./cmd/migrate -action create -dir migrations -name "$(name)"

.PHONY: run build test test-unit test-integration tidy swagger-setup swagger migrate-up migrate-down migrate-refresh migrate-create
