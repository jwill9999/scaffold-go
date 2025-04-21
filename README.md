# Go API Scaffolding System

A comprehensive scaffolding system for building production-ready Go API services with clean architecture, modern tooling, and best practices.

## Features

### Core Components
- ✅ Clean Architecture Implementation
- ✅ Modular Design
- ✅ Dependency Injection
- ✅ Configuration Management
- ✅ Error Handling
- ✅ Logging System

### Security
- ✅ CORS Middleware
- ✅ Rate Limiting (Redis and In-Memory)
- ✅ Security Headers
- ✅ JWT Authentication
- ❌ Input Validation
- ❌ Password Hashing
- ❌ CSRF Protection
- ❌ Session Management
- ❌ API Key Authentication
- ❌ OAuth2 Integration
- ❌ RBAC Implementation

### Database
- ✅ PostgreSQL Implementation
- ❌ MySQL Implementation
- ❌ MongoDB Implementation
- ❌ SQLite Implementation
- ❌ Migration System
- ❌ Connection Pooling
- ❌ Query Builders
- ❌ Transaction Management

### Caching
- ✅ Redis Implementation
- ✅ In-Memory Implementation
- ✅ Cache Interface
- ✅ TTL Management
- ✅ Atomic Operations

### Monitoring
- ✅ Prometheus Metrics
- ✅ Jaeger Tracing
- ❌ Custom Metric Collectors
- ❌ Alerting Rules
- ❌ Dashboard Templates
- ❌ Log Aggregation

### API Features
- ❌ Pagination
- ❌ Filtering
- ❌ Sorting
- ❌ Search
- ❌ Bulk Operations
- ❌ Webhook Support

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-scaffold.git

# Install dependencies
go mod download

# Build the scaffolding tool
make build
```

## Quick Start

```bash
# Initialize a new project
go-scaffold init my-project

# Generate components
go-scaffold generate model User
go-scaffold generate handler User
go-scaffold generate service User
go-scaffold generate repository User

# Run the project
make run
```

## Available Templates

### Core Components
- `model.go.tmpl` - Data models
- `handler.go.tmpl` - HTTP handlers
- `service.go.tmpl` - Business logic
- `repository.go.tmpl` - Data access
- `middleware.go.tmpl` - HTTP middleware
- `config.go.tmpl` - Configuration
- `errors.go.tmpl` - Error handling
- `logger.go.tmpl` - Logging

### Security
- `cors.go.tmpl` - CORS middleware
- `ratelimit.go.tmpl` - Rate limiting
- `headers.go.tmpl` - Security headers
- `jwt.go.tmpl` - JWT authentication

### Infrastructure
- `Dockerfile.tmpl` - Docker configuration
- `docker-compose.yml.tmpl` - Docker Compose setup
- `prometheus.yml.tmpl` - Prometheus configuration
- `jaeger.yml.tmpl` - Jaeger configuration

## Configuration

The scaffolding system uses a layered configuration approach:

```yaml
app:
  name: "my-api"
  version: "1.0.0"
  environment: "development"

server:
  port: 8080
  timeout: 30s
  cors:
    allowed_origins: ["*"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["Content-Type", "Authorization"]

database:
  host: "localhost"
  port: 5432
  name: "mydb"
  user: "user"
  password: "password"
  ssl_mode: "disable"

cache:
  type: "redis"
  host: "localhost"
  port: 6379
  password: ""
  db: 0

security:
  jwt_secret: "your-secret-key"
  token_expiry: "24h"
  rate_limit:
    requests_per_minute: 60
    burst_size: 10
```

## Development

```bash
# Build the tool
make build

# Run tests
make test

# Validate templates
make validate

# List available templates
make list-templates
```

## Roadmap

### Phase 1: Core Security (Current)
- [ ] Input Validation
- [ ] Password Hashing
- [ ] CSRF Protection
- [ ] Session Management

### Phase 2: Database & Authentication
- [ ] MySQL Implementation
- [ ] MongoDB Implementation
- [ ] OAuth2 Integration
- [ ] API Key Authentication

### Phase 3: API Features
- [ ] Pagination
- [ ] Filtering
- [ ] Sorting
- [ ] Search
- [ ] Bulk Operations

### Phase 4: Monitoring & Documentation
- [ ] Custom Metric Collectors
- [ ] Alerting Rules
- [ ] Dashboard Templates
- [ ] API Documentation

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