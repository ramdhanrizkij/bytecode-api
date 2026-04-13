-include .env
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export

# ─── App ─────────────────────────────────────────────────────────────────────

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

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

# ─── Migrations ──────────────────────────────────────────────────────────────
# Requires the `migrate` CLI: https://github.com/golang-migrate/migrate
# Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#
# DATABASE_URL can be overridden on the command line:
#   make migrate-up DATABASE_URL=postgres://...

DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Setup the migrate CLI tool
migrate-setup:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply all pending migrations
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose up

# Roll back the last applied migration
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose down 1

# Drop all migrations and apply them from scratch
migrate-refresh:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose down -all
	migrate -path migrations -database "$(DATABASE_URL)" -verbose up

# Create a new migration pair
# Usage: make migrate-create name=create_something_table
migrate-create:
	@test -n "$(name)" || (echo "ERROR: name is required. Usage: make migrate-create name=create_xxx_table" && exit 1)
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: run build test test-unit test-integration migrate-setup migrate-up migrate-down migrate-refresh migrate-create
