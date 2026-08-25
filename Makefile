.PHONY: run build test lint tidy \
        migrate-up migrate-down migrate-reset migrate-create \
        docker-up docker-down docker-all docker-migrate \
        sqlc-gen help

# ─── Config ───────────────────────────────────────────────────────────────────
APP_NAME  = mentorhub
MAIN_PATH = ./cmd/server
BIN_DIR   = ./bin

DB_USER     ?= mentorhub
DB_PASSWORD ?= mentorhub_secret
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_NAME     ?= mentorhub_db
DB_URL      = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ─── Development ──────────────────────────────────────────────────────────────

## run: запустить dev-сервер
run:
	go run $(MAIN_PATH)/main.go

## build: собрать бинарник
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-w -s" -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)/main.go
	@echo "✅ Built: $(BIN_DIR)/$(APP_NAME)"

## test: запустить тесты с покрытием
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

## lint: запустить golangci-lint
lint:
	golangci-lint run ./...

## tidy: привести зависимости в порядок
tidy:
	go mod tidy
	go mod verify

# ─── Migrations ───────────────────────────────────────────────────────────────

## migrate-up: применить все миграции
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

## migrate-down: откатить последнюю миграцию
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

## migrate-reset: сбросить все миграции (осторожно!)
migrate-reset:
	migrate -path migrations -database "$(DB_URL)" drop -f

## migrate-create name=<name>: создать новую миграцию
migrate-create:
	@[ "$(name)" ] || ( echo "❌ Usage: make migrate-create name=migration_name"; exit 1 )
	migrate create -ext sql -dir migrations -seq $(name)

# ─── Docker ───────────────────────────────────────────────────────────────────

## docker-up: поднять только PostgreSQL
docker-up:
	docker compose up -d mentorhub-db
	@echo "✅ PostgreSQL (mentorhub-db) запущен"

## docker-down: остановить все контейнеры
docker-down:
	docker compose down

## docker-all: поднять PostgreSQL + приложение
docker-all:
	docker compose up -d

## docker-migrate: запустить миграции через Docker
docker-migrate:
	docker compose --profile migrate run --rm mentorhub-migrate

## docker-logs: логи приложения
docker-logs:
	docker compose logs -f mentorhub-backend

# ─── Code Generation ──────────────────────────────────────────────────────────

## sqlc-gen: сгенерировать Go-код из SQL-запросов
sqlc-gen:
	sqlc generate

# ─── Help ─────────────────────────────────────────────────────────────────────

## help: показать список команд
help:
	@echo "MentorHub — доступные команды:"
	@grep -E '^## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ": "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' | \
		sed 's/## //'
