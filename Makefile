# Environment commands
COMPOSE_PROD=docker-compose -f docker/prod/docker-compose.prod.yml
COMPOSE_DEV=docker-compose -f docker/dev/docker-compose.yml

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=go-scaffold
VERSION?=0.0.1
BUILD_DIR=./bin
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
MAIN_FILE=./cmd/scaffold/main.go
CGO_ENABLED=0

# Project parameters (can be overridden by user)
PROJECT_NAME?=example-api
MODULE_PATH?=github.com/example/example-api
FEATURES?=
DB_TYPE?=postgres
DEPLOYMENT?=docker

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean install generate \
        test test-ci test-html test-auth test-register test-login test-coverage \
        lint lint-fix security scan-deps \
        prod-up prod-down prod-restart prod-logs prod-status prod-db-check prod-migrations-check prod-clean \
        dev-up dev-down dev-restart dev-logs dev-status dev-db-check dev-migrations-check dev-clean \
        example version release list-templates validate-templates run help

#----------------------------------------------
# Core Build Commands
#----------------------------------------------

## all: Clean, build, and test the project
all: clean build test

## build: Build the scaffold CLI tool
build:
	mkdir -p $(BUILD_DIR)
	$(CGO_ENABLED) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

## install: Install the scaffold CLI tool
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

## generate: Generate a new project with customizable settings
##           Usage: make generate PROJECT_NAME=myapi MODULE_PATH=github.com/username/myapi FEATURES=auth,metrics DB_TYPE=postgres
generate: build
	$(BUILD_DIR)/$(BINARY_NAME) -name $(PROJECT_NAME) -module $(MODULE_PATH) $(if $(FEATURES),-features $(FEATURES),) $(if $(DB_TYPE),-db $(DB_TYPE),) $(if $(DEPLOYMENT),-deployment $(DEPLOYMENT),)

## tidy: Tidy go modules
tidy:
	$(GOCMD) mod tidy

## update: Update dependencies
update:
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy

#----------------------------------------------
# Testing Commands
#----------------------------------------------

## test: Run all tests
test:
	@echo "Running all tests..."
	$(GOTEST) -v -race -cover ./tools/scaffold/...

## test-ci: Run tests with CI coverage (for GitHub Actions)
test-ci:
	@echo "Running tests with CI coverage..."
	$(GOTEST) -coverprofile=cover.out -json ./... | tee test-report.json

## test-html: Generate HTML coverage report
test-html:
	@echo "Generating HTML coverage report..."
	$(GOCMD) tool cover -html=cover.out -o coverage.html

## test-auth: Run auth handler tests
test-auth:
	@echo "Running auth handler tests..."
	$(GOTEST) -v ./handlers

## test-register: Run only registration tests
test-register:
	@echo "Running registration tests..."
	$(GOTEST) -v -run TestAuthHandler_Register ./handlers

## test-login: Run only login tests
test-login:
	@echo "Running login tests..."
	$(GOTEST) -v -run TestAuthHandler_Login ./handlers

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

#----------------------------------------------
# Code Quality Commands
#----------------------------------------------

## lint: Run linters
lint:
	golangci-lint run --no-config --timeout=5m ./...

## lint-fix: Run linters and fix issues
lint-fix:
	golangci-lint run --fix --no-config --timeout=5m ./...

## fmt: Format code
fmt:
	$(GOCMD) fmt ./...
	gofmt -s -w .
	goimports -w .

#----------------------------------------------
# Security Commands
#----------------------------------------------

## security: Run security checks with gosec
security:
	@echo "Running security checks..."
	gosec -fmt=json -out=security-report.json ./...
	@echo "Security report generated at security-report.json"

## scan-deps: Scan dependencies for vulnerabilities
scan-deps:
	@echo "Scanning dependencies for vulnerabilities..."
	go list -json -deps ./... | nancy sleuth

#----------------------------------------------
# Production Environment Commands
#----------------------------------------------

## prod-up: Start production containers
prod-up:
	$(COMPOSE_PROD) up -d

## prod-down: Stop production containers
prod-down:
	$(COMPOSE_PROD) down

## prod-clean: Stop containers and remove volumes (clean state)
prod-clean:
	$(COMPOSE_PROD) down -v

## prod-restart: Restart production containers
prod-restart: prod-down prod-up

## prod-logs: Show logs of all containers
prod-logs:
	$(COMPOSE_PROD) logs

## prod-logs-service: Show logs of specific container (usage: make prod-logs-service service=migrations)
prod-logs-service:
	$(COMPOSE_PROD) logs $(service)

## prod-status: Show status of containers
prod-status:
	$(COMPOSE_PROD) ps

## prod-db-check: Check database tables
prod-db-check:
	@echo "List of tables in the database:"
	@$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go -c "\dt" | cat
	@echo "\nDetailed table information:"
	@$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go -c "\d+ users" | cat

## prod-migrations-check: Check migration status
prod-migrations-check:
	@echo "Migration Status:"
	@$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go -c "SELECT * FROM schema_migrations;" | cat

## prod-db-shell: Run a psql shell
prod-db-shell:
	$(COMPOSE_PROD) exec postgres psql -U postgres -d scaffold-go

#----------------------------------------------
# Development Environment Commands
#----------------------------------------------

## dev-up: Start development containers
dev-up:
	$(COMPOSE_DEV) up -d

## dev-down: Stop development containers
dev-down:
	$(COMPOSE_DEV) down

## dev-clean: Stop containers and remove volumes
dev-clean:
	$(COMPOSE_DEV) down -v

## dev-restart: Restart development containers
dev-restart: dev-down dev-up

## dev-logs: Show logs of all containers
dev-logs:
	$(COMPOSE_DEV) logs

## dev-logs-service: Show logs of specific service (usage: make dev-logs-service service=migrations)
dev-logs-service:
	$(COMPOSE_DEV) logs $(service)

## dev-status: Show status of containers
dev-status:
	$(COMPOSE_DEV) ps

## dev-db-check: Check database tables
dev-db-check:
	@echo "List of tables in the database:"
	@$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go -c "\dt" | cat
	@echo "\nDetailed table information:"
	@$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go -c "\d+ users" | cat

## dev-migrations-check: Check migration status
dev-migrations-check:
	@echo "Migration Status:"
	@$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go -c "SELECT * FROM schema_migrations;" | cat

## dev-db-shell: Open a PostgreSQL shell
dev-db-shell:
	$(COMPOSE_DEV) exec postgres psql -U postgres -d scaffold-go

#----------------------------------------------
# Utility Commands
#----------------------------------------------

## example: Show example commands
example:
	@echo "Example commands:"
	@echo "  Create new project:"
	@echo "    $(BINARY_NAME) -name myapi -module github.com/username/myapi -features auth,metrics"
	@echo ""
	@echo "  Run the generated project:"
	@echo "    cd myapi && go mod tidy && go run ./cmd/api"

## version: Show version
version:
	@echo $(VERSION)

## release: Create a new release
release:
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	git push origin v$(VERSION)

## list-templates: List all available templates
list-templates:
	@find tools/scaffold/templates -name "*.tmpl" -type f | sort

## validate-templates: Validate all templates
validate-templates:
	@echo "Validating templates..."
	@find tools/scaffold/templates -name "*.tmpl" -type f -exec sh -c 'echo "Validating {}..." && cat {} > /dev/null' \;

## run: Build and run the example project
##      Usage: make run PROJECT_NAME=myapi MODULE_PATH=github.com/username/myapi
run: generate
	cd $(PROJECT_NAME) && $(GOCMD) mod tidy && $(GOCMD) run ./cmd/api

## Help:
.DEFAULT_GOAL := help
help:
	@echo "Go API Scaffolding System - Make Commands"
	@echo ""
	@echo "Core Build Commands:"
	@echo "  make all                  - Clean, build, and test the project"
	@echo "  make build                - Build the scaffold CLI tool"
	@echo "  make clean                - Remove build artifacts"
	@echo "  make install              - Install the scaffold CLI tool"
	@echo "  make generate             - Generate a new project (see customization options below)"
	@echo "  make tidy                 - Tidy go modules"
	@echo "  make update               - Update dependencies"
	@echo ""
	@echo "Project Customization Options:"
	@echo "  PROJECT_NAME=myapi        - Set the name of the generated project"
	@echo "  MODULE_PATH=github.com/user/repo - Set the Go module path"
	@echo "  FEATURES=auth,metrics     - Comma-separated list of features to include"
	@echo "  DB_TYPE=postgres          - Database type (postgres, mysql)"
	@echo "  DEPLOYMENT=docker         - Deployment type (docker, kubernetes)"
	@echo ""
	@echo "  Example: make generate PROJECT_NAME=my-service MODULE_PATH=github.com/myorg/my-service FEATURES=auth,metrics"
	@echo ""
	@echo "Testing Commands:"
	@echo "  make test                 - Run all tests"
	@echo "  make test-ci              - Run tests with CI coverage"
	@echo "  make test-html            - Generate HTML coverage report"
	@echo "  make test-auth            - Run auth handler tests"
	@echo "  make test-register        - Run only registration tests"
	@echo "  make test-login           - Run only login tests"
	@echo "  make test-coverage        - Run tests with coverage"
	@echo ""
	@echo "Code Quality Commands:"
	@echo "  make lint                 - Run linters"
	@echo "  make lint-fix             - Run linters and fix issues"
	@echo "  make fmt                  - Format code"
	@echo ""
	@echo "Security Commands:"
	@echo "  make security             - Run security checks with gosec"
	@echo "  make scan-deps            - Scan dependencies for vulnerabilities"
	@echo ""
	@echo "Production Environment Commands:"
	@echo "  make prod-up              - Start production containers"
	@echo "  make prod-down            - Stop production containers"
	@echo "  make prod-clean           - Stop containers and remove volumes"
	@echo "  make prod-restart         - Restart production containers"
	@echo "  make prod-logs            - Show logs of all containers"
	@echo "  make prod-logs-service    - Show logs of specific container"
	@echo "  make prod-status          - Show status of containers"
	@echo "  make prod-db-check        - Check database tables"
	@echo "  make prod-migrations-check - Check migration status"
	@echo "  make prod-db-shell        - Run a psql shell"
	@echo ""
	@echo "Development Environment Commands:"
	@echo "  make dev-up               - Start development containers"
	@echo "  make dev-down             - Stop development containers"
	@echo "  make dev-clean            - Stop containers and remove volumes"
	@echo "  make dev-restart          - Restart development containers"
	@echo "  make dev-logs             - Show logs of all containers"
	@echo "  make dev-logs-service     - Show logs of specific service"
	@echo "  make dev-status           - Show status of containers"
	@echo "  make dev-db-check         - Check database tables"
	@echo "  make dev-migrations-check - Check migration status"
	@echo "  make dev-db-shell         - Open a PostgreSQL shell"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make example              - Show example commands"
	@echo "  make version              - Show version"
	@echo "  make release              - Create a new release"
	@echo "  make list-templates       - List all available templates"
	@echo "  make validate-templates   - Validate all templates"
	@echo "  make run                  - Build and run the example project"
	@echo "  make help                 - Show this help message" 