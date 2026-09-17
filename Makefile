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
 
GREEN := \033[32m 
YELLOW := \033[33m 
BLUE := \033[34m 
RED := \033[31m 
RESET := \033[0m 
 
.PHONY: all help setup doctor run run-worker build build-api build-worker \
	fmt vet lint test test-race tidy deps \
	migrate-up migrate-down migrate-version migrate-create \
	docker-up docker-down docker-build docker-rebuild docker-logs docker-ps \
	swagger check clean 
 
all: check 
 
help: 
	@printf "$(BLUE)$(APP_NAME)$(RESET)\n" 
	@printf "\n" 
	@printf "$(GREEN)Usage:$(RESET)\n" 
	@printf "  make $(YELLOW)<target>$(RESET)\n" 
	@printf "\n" 
	@printf "$(GREEN)Development:$(RESET)\n" 
	@printf "  $(YELLOW)run$(RESET)              Run API server\n" 
	@printf "  $(YELLOW)run-worker$(RESET)       Run monitoring worker\n" 
	@printf "  $(YELLOW)build$(RESET)            Build API and worker\n" 
	@printf "  $(YELLOW)fmt$(RESET)              Format source code\n" 
	@printf "  $(YELLOW)vet$(RESET)              Run go vet\n" 
	@printf "  $(YELLOW)lint$(RESET)             Run golangci-lint\n" 
	@printf "  $(YELLOW)test$(RESET)             Run tests\n" 
	@printf "  $(YELLOW)test-race$(RESET)        Run tests with race detector\n" 
	@printf "\n" 
	@printf "$(GREEN)Database:$(RESET)\n" 
	@printf "  $(YELLOW)migrate-up$(RESET)       Apply migrations\n" 
	@printf "  $(YELLOW)migrate-down$(RESET)     Rollback last migration\n" 
	@printf "  $(YELLOW)migrate-version$(RESET)  Show migration version\n" 
	@printf "  $(YELLOW)migrate-create$(RESET)   Create migration\n" 
	@printf "\n" 
	@printf "$(GREEN)Docker:$(RESET)\n" 
	@printf "  $(YELLOW)docker-up$(RESET)         Start containers\n" 
	@printf "  $(YELLOW)docker-down$(RESET)       Stop containers\n" 
	@printf "  $(YELLOW)docker-build$(RESET)      Build containers\n" 
	@printf "  $(YELLOW)docker-rebuild$(RESET)    Rebuild containers\n" 
	@printf "  $(YELLOW)docker-logs$(RESET)       Show container logs\n" 
	@printf "  $(YELLOW)docker-ps$(RESET)         Show containers\n" 
	@printf "\n" 
	@printf "$(GREEN)Other:$(RESET)\n" 
	@printf "  $(YELLOW)setup$(RESET)             Check development environment\n" 
	@printf "  $(YELLOW)doctor$(RESET)            Check required tools\n" 
	@printf "  $(YELLOW)swagger$(RESET)           Generate Swagger docs\n" 
	@printf "  $(YELLOW)check$(RESET)             Run project checks\n" 
	@printf "  $(YELLOW)clean$(RESET)             Remove build artifacts\n" 
 
setup: 
	@$(MAKE) doctor 
	@$(GO) mod download 
 
doctor: 
	@printf "$(BLUE)Checking development environment...$(RESET)\n" 
	@command -v $(GO) >/dev/null 2>&1 || \
		(echo "$(RED)Go not found$(RESET)" && exit 1) 
	@printf "$(GREEN)Go:$(RESET) "; $(GO) version 
	@command -v $(DOCKER) >/dev/null 2>&1 || \
		(echo "$(RED)Docker not found$(RESET)" && exit 1) 
	@printf "$(GREEN)Docker:$(RESET) "; $(DOCKER) --version 
	@command -v $(GOOSE) >/dev/null 2>&1 || \
		(echo "$(YELLOW)goose not found$(RESET)")
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || \
		(echo "$(YELLOW)golangci-lint not found$(RESET)") 
 
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
	@if "$(name)"=="" (echo Usage: make migrate-create name=create_users && exit /b 1)
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