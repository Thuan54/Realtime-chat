# Makefile
.PHONY: help setup dev migrate test down

COMPOSE_FILE := infra/docker-compose.yml

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

setup: ## Bootstrap environment, validate .env, install dependencies
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env from .env.example. Please review values before starting."; fi
	@echo "Installing Go live-reload tooling (air)..."
	go install github.com/cosmtrek/air@latest
	@echo "Installing Frontend dependencies..."
	cd frontend && npm ci
	@echo "Pulling required Docker images..."
	docker compose -f $(COMPOSE_FILE) pull postgres redis

dev: ## Start local dev with concurrent Vite HMR + Go Air
	@echo "Starting infrastructure (Postgres, Redis)..."
	docker compose -f $(COMPOSE_FILE) up -d postgres redis
	@echo "Waiting for databases to initialize..."
	@sleep 5
	@echo "Launching Frontend (Vite) & Backend (Air) concurrently..."
	@cd frontend && npx vite & \
	 cd backend && air & \
	 wait

DB_URL := postgres://${POSTGRES_USER:-chatuser}:${POSTGRES_PASSWORD:-changeme}@localhost:5432/${POSTGRES_DB:-chatdb}?sslmode=disable

migrate-up: ## Apply all pending goose migrations
	@echo "Running goose up..."
	@cd migrations && goose -dir . postgres "$(DB_URL)" up

migrate-down: ## Rollback the last applied migration
	@echo "Running goose down..."
	@cd migrations && goose -dir . postgres "$(DB_URL)" down

migrate-status: ## Show migration version & status
	@cd migrations && goose -dir . postgres "$(DB_URL)" status

migrate-create: ## Create new SQL migration (usage: make migrate-create NAME=create_users)
ifndef NAME
	$(error NAME is required. Usage: make migrate-create NAME=create_channels)
endif
	@cd migrations && goose -dir . create $(NAME) sql

# --- QA / Testing Targets ---

test: test-unit test-integration ## Run all unit + integration tests

test-unit: ## Run frontend + backend unit tests only
	@echo "Running frontend unit tests..."
	cd frontend && npx vitest run --coverage
	@echo "Running backend unit tests..."
	cd backend && go test -v -race ./internal/... ./cmd/...

test-integration: ## Run backend integration tests with transactional isolation
	@echo "Running backend integration tests..."
	cd backend && go test -v -race -tags=integration ./internal/...

test-e2e: ## Run Playwright E2E against containerized staging env
	@echo "Starting E2E environment..."
	docker compose -f infra/docker-compose.yml up -d postgres redis
	@make migrate-up seed
	docker compose -f infra/docker-compose.yml up -d backend frontend nginx
	@echo "Running Playwright tests..."
	cd frontend && npx playwright test
	@echo "Tearing down..."
	docker compose -f infra/docker-compose.yml down -v

test-coverage: ## Generate combined coverage report (placeholder for CI aggregation)
	@echo "Coverage aggregation requires CI artifact merging; see 10-cicd-workflow.md"

down: ## Stop containers, preserve volumes
	@echo "Stopping development environment..."
	docker compose -f $(COMPOSE_FILE) down