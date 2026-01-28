# Google Lens AI Agent Makefile

# Variables
APP_NAME = google-lens-agent
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Docker variables
DOCKER_IMAGE = $(APP_NAME):$(VERSION)
DOCKER_REGISTRY = your-registry.com
DOCKER_REPO = $(DOCKER_REGISTRY)/$(APP_NAME)

# Go variables
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 1

.PHONY: help build test clean docker-build docker-run docker-push deploy dev test

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build targets
build: ## Build the application for current platform
	@echo "Building $(APP_NAME) v$(VERSION)..."
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

build-all: ## Build for all supported platforms
	@echo "Building for all platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(MAKE) build-rename
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 $(MAKE) build-rename
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(MAKE) build-rename
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 $(MAKE) build-rename

build-rename:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(APP_NAME)-$(GOOS)-$(GOARCH)$(if $(filter windows,$(GOOS)),.exe,) ./cmd/$(APP_NAME)

# Development targets
dev: ## Run in development mode
	@echo "Starting development server..."
	go run $(LDFLAGS) ./cmd/$(APP_NAME) serve -mode=interactive -verbose

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-coverage: test ## Open test coverage in browser
	@echo "Opening test coverage..."
	@open coverage.html || xdg-open coverage.html || echo "Coverage report saved to coverage.html"

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	@docker build -t $(DOCKER_IMAGE) .
	@docker tag $(DOCKER_IMAGE) $(APP_NAME):latest

docker-build-all: ## Build Docker image for all platforms
	@echo "Building Docker images for all platforms..."
	@docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_REPO):$(VERSION) -t $(DOCKER_REPO):latest --push .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 -v $(PWD)/data:/app/data $(DOCKER_IMAGE)

docker-compose-up: ## Start services with Docker Compose
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-compose-down: ## Stop services with Docker Compose
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

docker-compose-logs: ## Show Docker Compose logs
	@docker-compose logs -f

docker-push: docker-build ## Push Docker image to registry
	@echo "Pushing Docker image..."
	@docker tag $(DOCKER_IMAGE) $(DOCKER_REPO):$(VERSION)
	@docker push $(DOCKER_REPO):$(VERSION)
	@docker push $(DOCKER_REPO):latest

# Deployment targets
deploy-staging: ## Deploy to staging environment
	@echo "Deploying to staging..."
	@docker-compose -f docker-compose.yml -f docker-compose.staging.yml up -d

deploy-production: ## Deploy to production environment
	@echo "Deploying to production..."
	@docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Database targets
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@docker-compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f /app/migrations/001_init.sql

db-backup: ## Backup database
	@echo "Backing up database..."
	@docker-compose exec postgres pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB) > backup_$(shell date +%Y%m%d_%H%M%S).sql

db-restore: ## Restore database (usage: make db-restore BACKUP_FILE=backup.sql)
	@echo "Restoring database from $(BACKUP_FILE)..."
	@docker-compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < $(BACKUP_FILE)

# Utility targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@docker system prune -f

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

format: ## Format Go code
	@echo "Formatting Go code..."
	@go fmt ./...
	@goimports -w .

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

security-scan: ## Run security scan
	@echo "Running security scan..."
	@gosec ./...

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Monitoring targets
monitor: ## Open monitoring dashboards
	@echo "Opening monitoring dashboards..."
	@open http://localhost:3000 || xdg-open http://localhost:3000 || echo "Grafana: http://localhost:3000"
	@open http://localhost:9090 || xdg-open http://localhost:9090 || echo "Prometheus: http://localhost:9090"

logs: ## Show application logs
	@echo "Showing application logs..."
	@docker-compose logs -f $(APP_NAME)

# Setup targets
setup: ## Set up development environment
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

setup-ollama: ## Set up Ollama models
	@echo "Setting up Ollama models..."
	@curl -s http://localhost:11434/api/pull -X POST -d '{"name": "llava"}'
	@curl -s http://localhost:11434/api/pull -X POST -d '{"name": "llama2"}'

# CI/CD targets
ci: ## Run CI pipeline
	@echo "Running CI pipeline..."
	$(MAKE) deps
	$(MAKE) test
	$(MAKE) lint
	$(MAKE) security-scan
	$(MAKE) docker-build

release: ## Create a new release
	@echo "Creating release $(VERSION)..."
	@$(MAKE) test
	@$(MAKE) build-all
	@$(MAKE) docker-build
	@$(MAKE) docker-push
	@echo "Release $(VERSION) completed successfully"

# Version targets
version: ## Show version information
	@echo "$(APP_NAME) v$(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

info: ## Show build and runtime information
	@echo "Build Information:"
	@echo "  Version: $(VERSION)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(shell go version)"
	@echo "  Platform: $(shell go env GOOS)/$(shell go env GOARCH)"
	@echo ""
	@echo "Runtime Information:"
	@echo "  Docker: $(shell docker --version 2>/dev/null || echo 'Not installed')"
	@echo "  Docker Compose: $(shell docker-compose --version 2>/dev/null || echo 'Not installed')"
	@echo "  Ollama: $(shell curl -s http://localhost:11434/api/tags >/dev/null 2>&1 && echo 'Running' || echo 'Not running')"

# Quick commands
quick-start: setup docker-compose-up ## Quick setup and start
quick-dev: deps build dev ## Quick development setup
quick-test: build test ## Quick build and test

# Include custom makefile if it exists
-include Makefile.local