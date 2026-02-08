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

test: ## Run unit tests locally
	go test -v ./...

test-docker: ## Run unit tests inside a Docker container
	docker run --rm -v $$(pwd):/app -w /app golang:1.23-alpine go test -v ./...

tidy: ## Clean up go.mod and go.sum inside Docker
	docker run --rm -v $$(pwd):/app -w /app -e GOPROXY=https://goproxy.io,direct golang:1.23-alpine go mod tidy

up: ## Start the application in Docker (Detached)
	mkdir -p uploads
	docker build --network=host -t airgo:latest .
	docker compose up -d

rebuild: ## Force a clean rebuild without using cache
	mkdir -p uploads
	docker build --network=host --no-cache -t airgo:latest .
	docker compose up -d

prune: ## Remove all unused Docker build cache
	docker builder prune -f

down: ## Stop the Docker application
	docker compose down

logs: ## View Docker logs
	docker compose logs -f

url: ## Show the public Cloudflare Tunnel URL
	@docker compose logs tunnel 2>&1 | grep -o 'https://.*\.trycloudflare.com' | tail -n 1

qr: ## Show QR code for mobile access (requires qrencode)
	@IP=$$(hostname -I | awk '{print $$1}'); \
	URL="http://$$IP:8081"; \
	echo "\033[32m\033[1mAirGo LAN URL:\033[0m $$URL"; \
	if command -v qrencode >/dev/null 2>&1; then \
		qrencode -t ansiutf8 "$$URL"; \
	else \
		echo "\033[33mTip: Install 'qrencode' to show a QR code here (sudo apt install qrencode)\033[0m"; \
	fi

clean: ## Clean up binaries and temporary files
	rm -f $(APP_NAME)
	rm -rf uploads/*

security: ## Run basic security check (requires trivy installed, optional)
	@echo "Scanning for vulnerabilities..."
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image $(IMAGE_NAME)
