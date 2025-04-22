# Environment commands
COMPOSE_PROD=docker-compose -f docker/prod/docker-compose.prod.yml
COMPOSE_DEV=docker-compose -f docker/dev/docker-compose.yml

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=go-scaffold
VERSION?=0.0.1
BUILD_DIR=bin
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: prod-up prod-down prod-restart prod-logs prod-status prod-db-check prod-migrations-check prod-clean \
        dev-up dev-down dev-restart dev-logs dev-status dev-db-check dev-migrations-check dev-clean help \
        test test-auth test-register test-login test-coverage all build clean install generate lint fmt tidy update example version release list-templates validate-templates run

# Production environment commands
prod-up:
	$(COMPOSE_PROD) up -d

# Stop production containers
prod-down:
	$(COMPOSE_PROD) down

# Stop containers and remove volumes (clean state)
prod-clean:
	$(COMPOSE_PROD) down -v

# Restart production containers
prod-restart: prod-down prod-up

# Show logs of all containers
prod-logs:
	$(COMPOSE_PROD) logs

# Show logs of specific container (usage: make prod-logs-service service=migrations)
prod-logs-service:
	$(COMPOSE_PROD) logs $(service)

# Show status of containers
prod-status:
	$(COMPOSE_PROD) ps

# Check database tables
prod-db-check:
	@echo "List of tables in the database:"
	@$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go -c "\dt" | cat
	@echo "\nDetailed table information:"
	@$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go -c "\d+ users" | cat

# Check migration status
prod-migrations-check:
	@echo "Migration Status:"
	@$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go -c "SELECT * FROM schema_migrations;" | cat

# Run a psql shell
prod-db-shell:
	$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go

# Development environment commands
dev-up:
	$(COMPOSE_DEV) up -d

dev-down:
	$(COMPOSE_DEV) down

dev-clean:
	$(COMPOSE_DEV) down -v

dev-restart: dev-down dev-up

dev-logs:
	$(COMPOSE_DEV) logs

dev-logs-service:
	$(COMPOSE_DEV) logs $(service)

dev-status:
	$(COMPOSE_DEV) ps

dev-db-check:
	@echo "List of tables in the database:"
	@$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go -c "\dt" | cat
	@echo "\nDetailed table information:"
	@$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go -c "\d+ users" | cat

dev-migrations-check:
	@echo "Migration Status:"
	@$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go -c "SELECT * FROM schema_migrations;" | cat

dev-db-shell:
	$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go

# Test commands

# Run all tests
test:
	@echo "Running all tests..."
	$(GOTEST) -v -race -cover ./tools/scaffold/...

# Run auth handler tests
test-auth:
	@echo "Running auth handler tests..."
	$(GOTEST) -v ./handlers

# Run only registration tests
test-register:
	@echo "Running registration tests..."
	$(GOTEST) -v -run TestAuthHandler_Register ./handlers

# Run only login tests
test-login:
	@echo "Running login tests..."
	$(GOTEST) -v -run TestAuthHandler_Login ./handlers

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Help command
help:
	@echo "Production Environment Commands:"
	@echo "  make prod-up              - Start production containers"
	@echo "  make prod-down            - Stop production containers"
	@echo "  make prod-clean           - Stop containers and remove volumes"
	@echo "  make prod-restart         - Restart production containers"
	@echo "  make prod-logs            - Show logs of all containers"
	@echo "  make prod-logs-service    - Show logs of specific service (usage: make prod-logs-service service=migrations)"
	@echo "  make prod-status          - Show status of containers"
	@echo "  make prod-db-check        - Check database tables"
	@echo "  make prod-migrations-check - Check migration status"
	@echo "  make prod-db-shell        - Open a PostgreSQL shell"
	@echo "\nDevelopment Environment Commands:"
	@echo "  make dev-up               - Start development containers"
	@echo "  make dev-down             - Stop development containers"
	@echo "  make dev-clean            - Stop containers and remove volumes"
	@echo "  make dev-restart          - Restart development containers"
	@echo "  make dev-logs             - Show logs of all containers"
	@echo "  make dev-logs-service     - Show logs of specific service (usage: make dev-logs-service service=migrations)"
	@echo "  make dev-status           - Show status of containers"
	@echo "  make dev-db-check         - Check database tables"
	@echo "  make dev-migrations-check  - Check migration status"
	@echo "  make dev-db-shell         - Open a PostgreSQL shell"
	@echo "\nOther Commands:"
	@echo "  make all                  - Clean, build, and test the project"
	@echo "  make build                - Build the scaffold CLI tool"
	@echo "  make clean                - Remove build artifacts"
	@echo "  make install              - Install the scaffold CLI tool"
	@echo "  make generate             - Generate mock files for testing"
	@echo "  make lint                 - Run linters"
	@echo "  make fmt                  - Format code"
	@echo "  make tidy                 - Tidy go modules"
	@echo "  make update               - Update dependencies"
	@echo "  make example              - Show example commands"
	@echo "  make version              - Show version"
	@echo "  make release              - Create a new release"
	@echo "  make list-templates       - List all available templates"
	@echo "  make validate-templates   - Validate all templates"
	@echo "  make run                  - Build and run the example project"

# Build commands
all: clean build test

## Build:
build: ## Build the scaffold CLI tool
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) tools/scaffold/main.go

## Install:
install: build ## Install the scaffold CLI tool
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## Clean:
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

## Generate:
generate: build ## Generate a new example project with default settings
	$(BUILD_DIR)/$(BINARY_NAME) -name example-api -module github.com/example/example-api

## Lint:
lint: ## Run linters
	golangci-lint run

## Format:
fmt: ## Format code
	$(GOCMD) fmt ./...
	gofmt -s -w .

## Tidy:
tidy: ## Tidy go modules
	$(GOCMD) mod tidy

## Update:
update: ## Update dependencies
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy

## Example Commands:
example: ## Show example commands
	@echo "Example commands:"
	@echo "  Create new project:"
	@echo "    $(BINARY_NAME) -name myapi -module github.com/username/myapi"
	@echo ""
	@echo "  Run the generated project:"
	@echo "    cd myapi && go mod tidy && go run ./cmd/api"

## Version:
version: ## Show version
	@echo $(VERSION)

## Release:
release: ## Create a new release
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	git push origin v$(VERSION)

## Templates:
list-templates: ## List all available templates
	@find tools/scaffold/templates -name "*.tmpl" -type f | sort

## Validate:
validate-templates: ## Validate all templates
	@echo "Validating templates..."
	@find tools/scaffold/templates -name "*.tmpl" -type f -exec sh -c 'echo "Validating {}..." && cat {} > /dev/null' \;

## Run:
run: generate ## Build and run the example project
	cd example-api && $(GOCMD) mod tidy && $(GOCMD) run ./cmd/api

## Help:
.DEFAULT_GOAL := help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' 