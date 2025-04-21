# Go API Scaffolding System

This document describes how to use the Go API scaffolding system to generate and customize your API project.

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
--auth <type>        # Authentication type (jwt, oauth2, session)
--auth-providers     # OAuth2 providers (google, github, facebook)

# Database
--database <type>    # Database type (postgres, mysql, mongodb, sqlite)
--auto-migrate       # Enable auto-migrations in development

# Caching
--cache <type>       # Cache type (redis, memory)

# Observability
--metrics            # Enable Prometheus metrics
--tracing           # Enable distributed tracing
--logger <type>     # Logger type (zap, logrus)

# Documentation
--swagger            # Enable Swagger/OpenAPI documentation
--postman           # Generate Postman collection
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
- `string`: String field
- `int`: Integer field
- `float`: Float field
- `bool`: Boolean field
- `time`: Timestamp field
- `enum`: Enumerated type
- `uuid`: UUID field
- `json`: JSON field
- `array`: Array/slice field
- `map`: Map field
- `ref`: Reference to another model

## Project Structure

The scaffolding system generates the following structure:

```
myapi/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── models/              # Data models
│   ├── handlers/            # HTTP handlers
│   ├── services/            # Business logic
│   ├── repository/          # Data access
│   └── migrations/          # Database migrations
├── pkg/
│   ├── auth/               # Authentication utilities
│   ├── cache/              # Caching utilities
│   ├── database/           # Database utilities
│   ├── errors/             # Error handling
│   ├── logger/             # Logging utilities
│   ├── metrics/            # Metrics collection
│   └── tracing/            # Distributed tracing
├── scripts/                # Utility scripts
├── .air.toml              # Hot reload configuration
├── .gitignore
├── docker-compose.yml     # Development environment
├── Dockerfile             # Multi-stage build
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

# Add authentication
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

# Create Product resource with relationships
go-scaffold resource --name Product --fields "
  name:string:required
  description:string
  price:float:required,min=0
  user_id:ref:User
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

## Feature Details

### Authentication Options

1. **JWT Authentication**
```bash
go-scaffold init --name myapi --auth jwt
```
- Generates JWT middleware
- Token management
- Refresh token support
- Password hashing

2. **OAuth2 Authentication**
```bash
go-scaffold init --name myapi --auth oauth2 --auth-providers google,github
```
- OAuth2 provider integration
- Social login support
- Profile mapping

3. **Session Authentication**
```bash
go-scaffold init --name myapi --auth session
```
- Session management
- Cookie handling
- Session store (Redis/DB)

### Database Options

1. **PostgreSQL**
```bash
go-scaffold init --name myapi --database postgres
```
- GORM integration
- Migration support
- Connection pooling

2. **MongoDB**
```bash
go-scaffold init --name myapi --database mongodb
```
- MongoDB driver
- Document mapping
- Index management

### Caching Options

1. **Redis Cache**
```bash
go-scaffold init --name myapi --cache redis
```
- Redis client
- Cache middleware
- Distributed locking

2. **In-Memory Cache**
```bash
go-scaffold init --name myapi --cache memory
```
- Local caching
- TTL support
- Cache eviction

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