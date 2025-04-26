# Go API Scaffolding System

[![Go](https://github.com/jwill9999/scaffold-go/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/jwill9999/scaffold-go/actions/workflows/go.yml)
[![Release](https://github.com/jwill9999/scaffold-go/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/jwill9999/scaffold-go/actions/workflows/release.yml)
[![PR Check](https://github.com/jwill9999/scaffold-go/actions/workflows/pr-check.yml/badge.svg)](https://github.com/jwill9999/scaffold-go/actions/workflows/pr-check.yml)
[![Linting](https://github.com/jwill9999/scaffold-go/actions/workflows/linting.yml/badge.svg?branch=develop)](https://github.com/jwill9999/scaffold-go/actions/workflows/linting.yml)
[![Testing](https://github.com/jwill9999/scaffold-go/actions/workflows/testing.yml/badge.svg?branch=main)](https://github.com/jwill9999/scaffold-go/actions/workflows/testing.yml)
[![Security](https://github.com/jwill9999/scaffold-go/actions/workflows/security.yml/badge.svg?branch=main)](https://github.com/jwill9999/scaffold-go/actions/workflows/security.yml)
[![Quality](https://github.com/jwill9999/scaffold-go/actions/workflows/quality.yml/badge.svg?branch=main)](https://github.com/jwill9999/scaffold-go/actions/workflows/quality.yml)
[![codecov](https://codecov.io/gh/jwill9999/scaffold-go/branch/main/graph/badge.svg)](https://codecov.io/gh/jwill9999/scaffold-go)

A comprehensive scaffolding system for building production-ready Go API services with clean architecture, modern tooling, and best practices.

## Features

### Core Components

- ✅ Clean Architecture Implementation
- ✅ Modular Design
- ✅ Configuration Management (Viper)
- ✅ Structured Logging (Zap)
- ✅ HTTP Routing (Gin)
- ✅ Middleware Support
- ✅ Health Check Endpoints

### Database

- ✅ PostgreSQL Configuration

### Security

- ✅ Input Validation
- ✅ Security Headers Middleware
- ✅ Rate Limiting
- ✅ Secure Code Practices
- ✅ Security Scanning (gosec)
- ✅ Dependency Vulnerability Scanning

### Observability

- ✅ Structured Logging
- ✅ Request Tracing

## Installation

```bash
# Clone the repository
git clone https://github.com/jwill9999/scaffold-go.git

# Install dependencies
go mod download

# Build the scaffolding tool
make build
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/jwill9999/scaffold-go.git

# Install dependencies
go mod download

# Build the scaffolding tool
make build

# Generate a new project with default settings
make generate

# Generate a custom project with specific features
make generate PROJECT_NAME=my-api MODULE_PATH=github.com/username/my-api FEATURES=auth,metrics

# Or use the binary directly
./bin/go-scaffold -name my-api -module github.com/username/my-api -features auth,metrics

# Navigate to your project
cd my-api

# Initialize dependencies
go mod tidy

# Run the application
go run cmd/api/main.go
```

Your API will be available at `http://localhost:8080` with the following endpoints:

- Health check: `GET /health`
- API status: `GET /api/v1/status`

## Available Templates

The scaffolding system creates the following components:

### Project Structure

- `cmd/api/main.go` - Application entry point
- `internal/config/config.go` - Configuration management
- `internal/handlers/handlers.go` - HTTP handlers
- `internal/models/models.go` - Data models
- `internal/repository/repository.go` - Data access layer
- `internal/services/services.go` - Business logic
- `pkg/logger/logger.go` - Logging utilities
- `pkg/database/database.go` - Database connection
- `Dockerfile` - Container configuration
- `docker-compose.yml` - Multi-container setup

### Security Templates

- JWT Authentication
- Rate Limiting Middleware
- CORS Configuration
- Security Headers

## Development Commands

```bash
# Build the tool
make build

# Run tests
make test
npm run test:ci
npm run test:html

# Run linting
npm run lint

# Run security checks
npm run security

# Generate and run an example project
make run

# List available templates
make list-templates

# Validate templates
make validate-templates

# Show available commands
make help
```

## Configuration

The generated project uses a layered configuration approach with environment variables and YAML files.
Default configuration is loaded from:

1. Default values in code
2. Configuration file (`config.yaml` if present)
3. Environment variables (with `APP_` prefix)

Example configuration for a generated project:

```yaml
# config.yaml
server:
  port: 8080
  timeout: 30

security:
  cors:
    allowed_origins: ["*"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
  rate_limit:
    requests_per_minute: 60
    burst: 10

log_level: "info"
```

## Project Structure

The generated project follows a clean architecture pattern with the following layers:

```
my-project/
├── cmd/api/            # Application entrypoints
├── internal/           # Private application code
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   ├── core/           # Core business logic
│   │   ├── middleware/ # HTTP middleware
│   │   ├── errors/     # Error handling
│   │   └── server/     # Server configuration
│   └── services/       # Business logic
├── pkg/                # Public libraries
│   ├── database/       # Database utilities
│   ├── logger/         # Logging utilities
│   ├── security/       # Security utilities
│   └── metrics/        # Metrics collection
├── config.yaml         # Configuration file
├── Dockerfile          # Container definition
└── docker-compose.yml  # Container orchestration
```

## Security Features

The scaffold includes several security features:

1. **Static Code Analysis**: Uses gosec to detect security issues in code
2. **Dependency Scanning**: Checks dependencies for known vulnerabilities
3. **Path Traversal Prevention**: Safe file handling with validation
4. **Input Validation**: Request validation middleware
5. **Rate Limiting**: Configurable rate limiting with Redis or in-memory storage
6. **Security Headers**: CORS, CSP, and other security headers
7. **Authentication**: JWT authentication with secure configuration

These checks are integrated into the development workflow through:

- Pre-commit hooks (security scanning)
- CI/CD pipelines (automated security checks)
- Weekly scheduled security scans

## CI/CD Integration

The project uses GitHub Actions for continuous integration and delivery with a modular workflow structure:

### Core Workflows

- **quality.yml**: Orchestrates the linting and testing processes by calling reusable workflows
- **security.yml**: Runs security scans (gosec) and dependency checks (nancy)
- **linting.yml**: Performs code formatting and linting checks (reusable)
- **testing.yml**: Executes tests with code coverage reporting (reusable)
- **release.yml**: Handles the release process with semantic versioning
- **pr-check.yml**: Validates pull requests for commit message standards and merge conflicts
- **codeql.yml**: Analyzes code for security vulnerabilities

### Workflow Structure

- **Reusable Workflows**: Both linting and testing are implemented as reusable workflows
- **Modular Design**: Each workflow has a specific responsibility
- **Reduced Duplication**: Common tasks are defined once and reused across workflows
- **Scheduled Scans**: Security workflows run weekly to ensure ongoing protection

This approach ensures maintainability, standardization, and efficient updates across the CI/CD pipeline.

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes using our commit script:
   ```bash
   npm run commit
   ```
   This will guide you through creating a properly formatted commit message
5. Push to your branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

For more details, see [our commit guide](docs/COMMIT_GUIDE.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.

## Acknowledgments

- Clean Architecture by Robert C. Martin
- Go best practices and patterns
- Open-source community contributions
