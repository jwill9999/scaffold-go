# Go API Scaffolding System

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
# Generate a new project
./bin/go-scaffold -name my-api -module github.com/username/my-api

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

## Development Commands

```bash
# Build the tool
make build

# Run tests
make test

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
│   └── services/       # Business logic
├── pkg/                # Public libraries
│   ├── database/       # Database utilities
│   └── logger/         # Logging utilities
├── config.yaml         # Configuration file
├── Dockerfile          # Container definition
└── docker-compose.yml  # Container orchestration
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.

## Acknowledgments

- Clean Architecture by Robert C. Martin
- Go best practices and patterns
- Open-source community contributions