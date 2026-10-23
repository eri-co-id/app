.PHONY: help build run test lint clean

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the server"
	@echo "  make test         - Run all tests"
	@echo "  make lint         - Run linters (gofmt, go vet)"
	@echo "  make fmt          - Format code with gofmt"
	@echo "  make check-fast   - Quick checks (fmt, vet, gitleaks)"
	@echo "  make check-task   - Task checks (tests + coverage)"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make dev-up       - Start docker-compose"
	@echo "  make dev-down     - Stop docker-compose"

build:
	@echo "Building..."
	go build -o bin/server ./cmd/server

run: build
	@echo "Running server on port 8080..."
	PORT=8080 ./bin/server

test:
	@echo "Running tests..."
	go test -v -race ./...

lint:
	@echo "Running linters..."
	gofmt -l .
	go vet ./...

fmt:
	@echo "Formatting code..."
	gofmt -w .

check-fast:
	@echo "Running fast checks..."
	gofmt -l .
	go vet ./...
	gitleaks detect --redact --no-banner

check-task: check-fast test
	@echo "Task checks complete"

clean:
	@echo "Cleaning..."
	rm -rf bin/

dev-up:
	@echo "Starting docker services..."
	docker compose up -d

dev-down:
	@echo "Stopping docker services..."
	docker compose down

.DEFAULT_GOAL := help
