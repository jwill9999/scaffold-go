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
- ✅ Middleware System
- ✅ Migration System
- ✅ Testing Framework
- ✅ Mock Generation

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
- ✅ Migration System
- ✅ Connection Pooling
- ✅ Query Builders
- ✅ Transaction Management

### Caching
- ✅ Redis Implementation
- ✅ In-Memory Implementation
- ✅ Cache Interface
- ✅ TTL Management
- ✅ Atomic Operations
- ✅ Distributed Locking

### Monitoring
- ✅ Prometheus Metrics
- ✅ Jaeger Tracing
- ✅ Custom Metric Collectors
- ❌ Alerting Rules
- ✅ Dashboard Templates
- ❌ Log Aggregation

### API Features
- ✅ Swagger/OpenAPI Documentation
- ✅ Postman Collection Generation
- ❌ Pagination
- ❌ Filtering
- ❌ Sorting
- ❌ Search
- ❌ Bulk Operations
- ❌ Webhook Support

### Infrastructure
- ✅ Docker Configuration
- ✅ Docker Compose Setup
- ✅ Environment Configuration
- ✅ Hot Reload Support
- ✅ Multi-stage Builds

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

The scaffolding system uses a layered configuration approach with environment variables and YAML files. Here's an example configuration:

```yaml
# config.yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 60s

database:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "postgres"
    dbname: "app"
    sslmode: "disable"
    max_connections: 25
    max_idle_connections: 5
    connection_lifetime: "1h"

cache:
  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    pool_size: 10
    min_idle_conns: 2
    max_retries: 3
    dial_timeout: "5s"
    read_timeout: "3s"
    write_timeout: "3s"

  memory:
    cleanup_interval: "1m"
    default_ttl: "5m"

security:
  jwt:
    secret: "your-secret-key"
    expiration: "24h"
    refresh_expiration: "168h"  # 7 days

  cors:
    allowed_origins: ["*"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allowed_headers: ["Origin", "Content-Type", "Accept", "Authorization"]
    exposed_headers: ["Content-Length"]
    allow_credentials: true
    max_age: "12h"

  rate_limit:
    requests_per_minute: 60
    burst_size: 10
    key_prefix: "rate_limit:"
    block_duration: "5m"

monitoring:
  prometheus:
    enabled: true
    path: "/metrics"
    port: 9090

  jaeger:
    enabled: true
    host: "localhost"
    port: 6831
    service_name: "app"
    sampler_type: "const"
    sampler_param: 1

logging:
  level: "info"
  format: "json"
  output: "stdout"
  file:
    path: "logs/app.log"
    max_size: 100
    max_backups: 3
    max_age: 7
    compress: true
```

The configuration can be overridden using environment variables with the `APP_` prefix:

```bash
APP_SERVER_PORT=8080
APP_DATABASE_POSTGRES_HOST=localhost
APP_CACHE_REDIS_HOST=redis
APP_SECURITY_JWT_SECRET=your-secret-key
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