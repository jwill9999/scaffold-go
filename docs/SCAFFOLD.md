# Go API Scaffolding System

This document describes how to use the Go API scaffolding system to generate and customize your API project.

## Implementation Status

This documentation includes both implemented features and planned features. Features are marked as follows:
- ✅ Implemented and available for use
- 🔜 Planned for future implementation

## Quick Start

```bash
# Initialize a new project
go-scaffold init --name myapi --module github.com/username/myapi

# Add resources
go-scaffold resource --name User --fields "name:string:required email:string:required,email"

# Start development
docker-compose up -d
```

## Command Reference

### Project Initialization

```bash
go-scaffold init [flags]
```

#### Required Flags
- `--name`: Project name
- `--module`: Go module path (e.g., github.com/username/project)

#### Optional Flags
- `--description`: Project description
- `--author`: Project author
- `--license`: License type (default: MIT)

### Feature Flags

```bash
# Authentication
--auth <type>        # Authentication type (✅ jwt, 🔜 oauth2, 🔜 session)
🔜 --auth-providers  # OAuth2 providers (google, github, facebook)

# Database
--database <type>    # Database type (✅ postgres, 🔜 mysql, 🔜 mongodb, 🔜 sqlite)
--auto-migrate       # Enable auto-migrations in development

# Caching
--cache <type>       # Cache type (✅ redis, ✅ memory)

# Observability
--metrics            # Enable Prometheus metrics (✅)
--tracing            # Enable distributed tracing (✅ Jaeger)
--logger <type>      # Logger type (✅ structured logging)

# Documentation
--swagger            # Enable Swagger/OpenAPI documentation (✅)
🔜 --postman         # Generate Postman collection
```

### Resource Generation

```bash
go-scaffold resource [flags]
```

#### Required Flags
- `--name`: Resource name (e.g., User, Product)
- `--fields`: Field definitions

#### Field Definition Format
```
name:type:validation[,validation]
```

Examples:
```bash
# Simple field
--fields "name:string:required"

# Multiple fields
--fields "name:string:required email:string:required,email age:int:required,min=18"

# Complex fields
--fields "
  name:string:required
  email:string:required,email
  password:string:required,min=8
  role:enum:required,oneof=admin|user
  status:bool:required
  created_at:time
"
```

#### Supported Field Types
Currently implemented:
- ✅ `string`: String field
- ✅ `int`: Integer field
- ✅ `float`: Float field
- ✅ `bool`: Boolean field
- ✅ `time`: Timestamp field
- ✅ `enum`: Enumerated type

Planned for future releases:
- 🔜 `uuid`: UUID field
- 🔜 `json`: JSON field
- 🔜 `array`: Array/slice field
- 🔜 `map`: Map field
- 🔜 `ref`: Reference to another model

## Project Structure

The scaffolding system generates the following structure:

```
myapi/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point (✅)
├── internal/
│   ├── config/              # Configuration management (✅)
│   ├── models/              # Data models (✅)
│   ├── handlers/            # HTTP handlers (✅)
│   ├── services/            # Business logic (✅)
│   ├── repository/          # Data access (✅)
│   └── migrations/          # Database migrations (✅)
├── pkg/
│   ├── auth/               # Authentication utilities (✅ JWT only)
│   ├── cache/              # Caching utilities (✅ Redis, Memory)
│   ├── database/           # Database utilities (✅ PostgreSQL only)
│   ├── errors/             # Error handling (✅)
│   ├── logger/             # Logging utilities (✅)
│   ├── metrics/            # Metrics collection (✅ Prometheus)
│   ├── tracing/            # Distributed tracing (✅ Jaeger)
│   └── security/           # Security utilities (✅ CORS, Rate limiting, Security headers)
├── scripts/                # Utility scripts (🔜)
├── .gitignore
├── docker-compose.yml     # Development environment (✅)
├── Dockerfile             # Container build (✅)
├── go.mod
└── README.md
```

## Configuration

The system uses a layered configuration approach:

1. Default values in code
2. Configuration files (YAML)
3. Environment variables
4. Command-line flags

### Configuration File (config.yaml)

```yaml
app:
  name: myapi
  version: 1.0.0
  environment: development

server:
  port: 8080
  timeout: 30s
  cors_enabled: true
  cors_origins:
    - http://localhost:3000

database:
  driver: postgres
  host: localhost
  port: 5432
  name: myapi
  user: postgres
  password: postgres
  ssl_mode: disable
  auto_migrate: true  # Development only

auth:
  jwt_secret: your-secret-key
  jwt_expiration_time: 24h
  refresh_token_duration: 720h
  password_min_length: 8

features:
  enable_tracing: true
  enable_metrics: true
  enable_cache: true
  enable_rate_limit: true
  enable_swagger: true
```

## Development Workflow

### 1. Initialize Project
```bash
# Create new project
go-scaffold init --name myapi --module github.com/username/myapi

# Add authentication (currently only JWT is fully implemented)
go-scaffold init --name myapi --auth jwt

# Add database and cache
go-scaffold init --name myapi --database postgres --cache redis
```

### 2. Generate Resources
```bash
# Create User resource
go-scaffold resource --name User --fields "
  name:string:required
  email:string:required,email
  password:string:required,min=8
  role:enum:required,oneof=admin|user
"

# Create Product resource
go-scaffold resource --name Product --fields "
  name:string:required
  description:string
  price:float:required,min=0
"
```

### 3. Database Migrations
```bash
# Create migration
go-scaffold migration create add_users_table

# Run migrations
go run cmd/migrate/main.go up

# Rollback migrations
go run cmd/migrate/main.go down
```

### 4. Start Development
```bash
# Start dependencies
docker-compose up -d

# Run with hot reload
air

# Run tests
go test ./...
```

## Implemented Features

### Authentication
- ✅ JWT Authentication
  - Token generation and validation
  - Refresh token support
  - Expiration handling

### Database Support
- ✅ PostgreSQL
  - Connection management
  - Migration support
  - Repository pattern implementation

### Caching
- ✅ Redis Cache
  - Key-value operations
  - TTL support
  - Atomic operations
- ✅ In-Memory Cache
  - Fast local caching
  - Expiration management
  - Thread-safe operations

### Security
- ✅ CORS middleware
  - Configurable origins, methods, headers
  - Preflight handling
- ✅ Rate limiting
  - Redis-backed implementation
  - In-memory implementation
  - Configurable limits and windows
- ✅ Security headers
  - Content security policy
  - XSS protection
  - Frame options

### Observability
- ✅ Prometheus metrics
  - Request timing
  - Error counting
  - Custom metric support
- ✅ Jaeger tracing
  - Distributed tracing
  - Span management
  - Context propagation
- ✅ Structured logging
  - Log levels
  - Contextual information
  - Error details

### Infrastructure
- ✅ Dockerfile
  - Build configuration
  - Runtime environment
- ✅ Docker Compose
  - Development setup
  - Service definitions
- ✅ Environment configuration
  - Example .env file
  - Configuration loading

### Documentation
- ✅ Swagger/OpenAPI
  - API documentation
  - Schema definitions
  - Endpoint specifications

## Coming Soon

The following features are planned for future releases:

### Authentication Enhancements
- 🔜 OAuth2 Authentication
  - Multiple provider support
  - Social login integration
- 🔜 Session-based Authentication
  - Cookie management
  - Session storage
- 🔜 Multi-factor Authentication
  - TOTP implementation
  - Recovery codes

### Additional Database Support
- 🔜 MySQL Support
- 🔜 MongoDB Support
- 🔜 SQLite Support
- 🔜 Advanced query builders
- 🔜 Transaction management

### API Features
- 🔜 Pagination
- 🔜 Filtering
- 🔜 Sorting
- 🔜 Search
- 🔜 Webhooks

### Security Enhancements
- 🔜 Input validation
- 🔜 Password hashing
- 🔜 CSRF protection
- 🔜 Role-based access control

### DevOps Support
- 🔜 CI/CD templates
- 🔜 Kubernetes manifests
- 🔜 Production deployment configurations

### Testing Enhancements
- 🔜 Integration test templates
- 🔜 End-to-end test templates
- 🔜 Load testing templates
- 🔜 Test data generation

## Best Practices

1. **Version Control**
- Use semantic versioning
- Tag releases
- Maintain changelog

2. **Testing**
- Write unit tests
- Integration tests
- E2E tests
- Use test containers

3. **Documentation**
- Update API documentation
- Document configuration
- Maintain examples

4. **Security**
- Use environment variables
- Secure passwords
- Rate limiting
- Input validation

## Troubleshooting

Common issues and solutions:

1. **Database Connection**
```bash
# Check database status
docker-compose ps

# View database logs
docker-compose logs postgres
```

2. **Migration Issues**
```bash
# Check migration status
go run cmd/migrate/main.go status

# Force migration version
go run cmd/migrate/main.go force 20230101000000
```

3. **Build Errors**
```bash
# Clean go cache
go clean -cache

# Update dependencies
go mod tidy
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License 