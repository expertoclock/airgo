# Makefile
# Variables
APP_NAME := airgo
IMAGE_NAME := airgo:latest

.PHONY: help build run up down clean test lint

help: ## Show this help message
	@echo "AirGo will be available at http://localhost:8081"
	@echo "Usage: make [target]"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go binary locally
	go build -o $(APP_NAME) .

run: ## Run the app locally (requires Go installed)
	go run main.go

test: ## Run unit tests
	go test -v ./...

up: ## Start the application in Docker (Detached)
	mkdir -p uploads
	docker compose up -d --build

down: ## Stop the Docker application
	docker compose down

logs: ## View Docker logs
	docker compose logs -f

clean: ## Clean up binaries and temporary files
	rm -f $(APP_NAME)
	rm -rf uploads/*

security: ## Run basic security check (requires trivy installed, optional)
	@echo "Scanning for vulnerabilities..."
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image $(IMAGE_NAME)
