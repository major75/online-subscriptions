# Load .env file and export its env variables
include .env
export

GOOSE_DIR=./database/migrations
GOOSE_DBSTRING=$(DATABASE_URL)

export GOOSE_DBSTRING
export GOOSE_DRIVER=postgres

.PHONY: install build run clean docker-build docker-up docker-down docker-logs docker-rebuild swagger migrate-create migrate-up migrate-down migrate-status

BINARY_NAME=subscriptions
DOCKER_IMAGE=subscriptions-service:latest

install: ## Install dependencies
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

build: ## Build the application
	go build -a -o bin/$(BINARY_NAME) ./cmd/api

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run test
	go test ./...

clean: ## Remove build artifacts
	rm -rf bin/

docker-build: ## Build Docker image
	# docker build -t $(DOCKER_IMAGE) .
	docker build --progress=plain -t $(DOCKER_IMAGE) .

docker-up: ## Start Docker containers
	docker compose up -d

docker-down: ## Stop Docker containers
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

docker-rebuild: ## Rebuild and restart Docker containers
	docker compose down
	docker build --progress=plain --no-cache -t $(DOCKER_IMAGE) .
	docker compose up -d

swagger: ## Generate Swagger documentation
	@echo "Installing swag..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating swagger documentation..."
	$(shell go env GOPATH)/bin/swag init -g cmd/api/main.go -o docs

migrate-status:
	goose -dir $(GOOSE_DIR) status

migrate-create:
	goose -dir $(GOOSE_DIR) create $(name) sql

migrate-up:
	goose -dir $(GOOSE_DIR) up

migrate-down:
	goose -dir $(GOOSE_DIR) down 1