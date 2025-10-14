.PHONY: help deps build run dev dev-stop dev-logs docker-build docker-up docker-down docker-logs generate test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the binary
	go build -o bin/server cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

dev: ## Run in development mode with hot reload
	docker-compose -f docker-compose.dev.yml up -d

dev-stop: ## Stop development containers
	docker-compose -f docker-compose.dev.yml down

dev-logs: ## View development logs
	docker-compose -f docker-compose.dev.yml logs -f

docker-build: ## Build production Docker image
	docker build -t poc-collaborative-filter:latest .

docker-up: ## Start production containers
	docker-compose up -d

docker-down: ## Stop production containers
	docker-compose down

docker-logs: ## View production logs
	docker-compose logs -f

generate: ## Generate GraphQL code
	go run github.com/99designs/gqlgen generate

test: ## Run tests
	go test -v ./...

seed: ## Seed database with mock data
	@echo "Seeding database with mockdata.json..."
	@./seed-data.sh || echo "Error: Seeding failed. Make sure services are running (make dev)"

test-setup: dev seed ## Start services and seed database
	@echo "✅ Test environment ready!"
	@echo "GraphQL endpoint: http://localhost:8080/query"
	@echo "GraphQL Playground: http://localhost:8080"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf graph/generated/

.DEFAULT_GOAL := help
