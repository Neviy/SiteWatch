APP_NAME ?= sitewatch

GO ?= go
GOFMT ?= gofmt
GOOSE ?= goose
DOCKER ?= docker
DOCKER_COMPOSE ?= docker compose
GOLANGCI_LINT ?= golangci-lint

API_BINARY ?= bin/api
WORKER_BINARY ?= bin/worker

MIGRATIONS_DIR ?= migrations
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/sitewatch?sslmode=disable

CHECK_INTERVAL ?= 60s
HTTP_TIMEOUT ?= 5s
WORKER_COUNT ?= 4

.PHONY: all help setup doctor run run-worker build build-api build-worker fmt vet lint test test-race tidy deps migrate-up migrate-down migrate-version migrate-create docker-up docker-down docker-build docker-rebuild docker-logs docker-ps swagger check clean

all: check

help:
	@printf "$(APP_NAME)\n"
	@printf "\n"
	@printf "Usage:\n"
	@printf "  make <target>\n"
	@printf "\n"
	@printf "Development:\n"
	@printf "  make run                    Run API server\n"
	@printf "  make run-worker            Run monitoring worker\n"
	@printf "  make build                 Build API and worker\n"
	@printf "  make fmt                   Format source code\n"
	@printf "  make vet                   Run go vet\n"
	@printf "  make lint                  Run golangci-lint\n"
	@printf "  make test                  Run tests\n"
	@printf "  make test-race             Run tests with race detector\n"
	@printf "\n"
	@printf "Database:\n"
	@printf "  make migrate-up            Apply migrations\n"
	@printf "  make migrate-down          Rollback last migration\n"
	@printf "  make migrate-version       Show migration version\n"
	@printf "  make migrate-create name=create_users\n"
	@printf "\n"
	@printf "Docker:\n"
	@printf "  make docker-up             Start containers\n"
	@printf "  make docker-down           Stop containers\n"
	@printf "  make docker-build          Build containers\n"
	@printf "  make docker-rebuild        Rebuild containers\n"
	@printf "  make docker-logs           Show container logs\n"
	@printf "  make docker-ps             Show containers\n"
	@printf "\n"
	@printf "Other:\n"
	@printf "  make setup                 Download dependencies\n"
	@printf "  make doctor                Check required tools\n"
	@printf "  make swagger               Generate Swagger docs\n"
	@printf "  make check                 Run project checks\n"
	@printf "  make clean                 Remove build artifacts\n"

setup:
	@$(MAKE) doctor
	@$(GO) mod download

doctor:
	@printf "Checking development environment...\n"
	@command -v $(GO) >/dev/null 2>&1 || { printf "Go not found\n"; exit 1; }
	@printf "Go:\n"
	@$(GO) version

	@command -v $(DOCKER) >/dev/null 2>&1 || { printf "Docker not found\n"; exit 1; }
	@printf "Docker:\n"
	@$(DOCKER) --version

	@command -v $(GOOSE) >/dev/null 2>&1 || printf "WARNING: goose not found\n"
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || printf "WARNING: golangci-lint not found\n"

run:
	CHECK_INTERVAL=$(CHECK_INTERVAL) \
	HTTP_TIMEOUT=$(HTTP_TIMEOUT) \
	WORKER_COUNT=$(WORKER_COUNT) \
	$(GO) run ./cmd/api

run-worker:
	CHECK_INTERVAL=$(CHECK_INTERVAL) \
	HTTP_TIMEOUT=$(HTTP_TIMEOUT) \
	WORKER_COUNT=$(WORKER_COUNT) \
	$(GO) run ./cmd/worker

build: build-api build-worker

build-api:
	@mkdir -p bin
	$(GO) build -o $(API_BINARY) ./cmd/api

build-worker:
	@mkdir -p bin
	$(GO) build -o $(WORKER_BINARY) ./cmd/worker

fmt:
	$(GOFMT) -w .

vet:
	$(GO) vet ./...

lint:
	$(GOLANGCI_LINT) run

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

tidy:
	$(GO) mod tidy

deps:
	$(GO) mod download

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-version:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" version

migrate-create:
	@test -n "$(name)" || { \
		printf "Usage: make migrate-create name=create_users\n"; \
		exit 1; \
	}
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-build:
	$(DOCKER_COMPOSE) build

docker-rebuild:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-ps:
	$(DOCKER_COMPOSE) ps

swagger:
	swag init -g cmd/api/main.go

check: fmt vet test

clean:
	rm -rf bin