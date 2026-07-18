.PHONY: run build test test-coverage lint sqlc-generate \
	migrate-up migrate-down migrate-status migrate-create \
	db-up db-down db-reset

GOOSE_DRIVER := postgres
GOOSE_DBSTRING := $(DATABASE_URL)
MIGRATIONS_DIR := db/migrations

run:
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -o bin/api ./cmd/api

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

sqlc-generate:
	sqlc generate

migrate-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) up

migrate-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) down

migrate-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(MIGRATIONS_DIR) status

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d postgres
