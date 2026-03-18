ifneq (,$(wildcard .env))
include .env
export
endif

GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose@v3.24.1
GOOSE_INSTALL ?= go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
MIGRATIONS_DIR ?= migrations
DB_DSN ?= host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE) TimeZone=UTC
BUILD_DIR ?= bin

.PHONY: install-goose build build-all build-api build-worker clean migrate-up migrate-down migrate-reset migrate-status goose-create run-api run-worker

install-goose:
	$(GOOSE_INSTALL)

build: build-api build-worker

build-all: clean build

build-api:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/api ./cmd/api

build-worker:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/worker ./cmd/worker

clean:
	rm -rf $(BUILD_DIR)

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

goose-create:
	@if [ -z "$(name)" ]; then echo "usage: make goose-create name=create_something"; exit 1; fi
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker