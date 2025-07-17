# Variables
DOCKER_COMPOSE = docker-compose -f deployments/docker-compose.yml
DOCKER_COMPOSE_DEV = docker-compose -f deployments/docker-compose.dev.yml
GO = go
GOFLAGS = -v
SERVICES = url-service user-service analytics-service

.PHONY: all build test clean

# Default target
all: build

# Initialize project
init:
	@echo "Initializing project..."
	@$(GO) work sync
	@$(MAKE) tools
	@$(MAKE) proto-gen

# Install development tools
tools:
	@echo "Installing development tools..."
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@$(GO) install github.com/vektra/mockery/v2@latest

# Build all services
build:
	@echo "Building services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd services/$$service && $(GO) build $(GOFLAGS) -o ../../bin/$$service ./cmd/server; \
		cd ../..; \
	done

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test ./... -cover -race

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@$(DOCKER_COMPOSE_DEV) up -d
	@$(GO) test ./... -tags=integration -cover -race
	@$(DOCKER_COMPOSE_DEV) down

# Run linters
lint:
	@echo "Running linters..."
	@golangci-lint run ./...

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@mockery --all --output=mocks --dir=.

# Proto generation
proto-gen:
	@echo "Generating proto files..."
	@buf generate

# Start development environment
dev-up:
	@echo "Starting development environment..."
	@$(DOCKER_COMPOSE_DEV) up -d

# Stop development environment
dev-down:
	@echo "Stopping development environment..."
	@$(DOCKER_COMPOSE_DEV) down

# View logs
logs:
	@$(DOCKER_COMPOSE_DEV) logs -f

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf vendor/
	@$(GO) clean -cache

# Database migrations
migrate-up:
	@echo "Running migrations up..."
	@migrate -path services/url-service/migrations -database "postgresql://postgres:postgres@localhost:5432/url_service?sslmode=disable" up

migrate-down:
	@echo "Running migrations down..."
	@migrate -path services/url-service/migrations -database "postgresql://postgres:postgres@localhost:5432/url_service?sslmode=disable" down

# Help
help:
	@echo "Available targets:"
	@echo "  init              - Initialize project"
	@echo "  build             - Build all services"
	@echo "  test              - Run unit tests"
	@echo "  test-integration  - Run integration tests"
	@echo "  lint              - Run linters"
	@echo "  mocks             - Generate mocks"
	@echo "  proto-gen         - Generate proto files"
	@echo "  dev-up            - Start development environment"
	@echo "  dev-down          - Stop development environment"
	@echo "  logs              - View docker logs"
	@echo "  clean             - Clean build artifacts"
	@echo "  migrate-up        - Run database migrations up"
	@echo "  migrate-down      - Run database migrations down"