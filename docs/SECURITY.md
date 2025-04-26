# Security Guide

This document outlines the security features, tools, and practices implemented in our Go API scaffolding system.

## Security Features

### Code Security

- **Static Code Analysis**: The project uses `gosec` to detect security issues like:

  - Command injection vulnerabilities
  - SQL injection risks
  - Path traversal vulnerabilities
  - Insecure file operations
  - Hard-coded credentials
  - Insecure cryptographic practices

- **Path Traversal Prevention**: The scaffolding tool implements comprehensive validation to prevent path traversal attacks:

  - Validation of input paths
  - Canonicalization of paths
  - Prevention of directory traversal through relative paths
  - Safe file operations with proper permissions

- **Command Injection Prevention**: The tool implements secure command execution:
  - Whitelisting of allowed commands
  - Input validation and sanitization
  - Use of proper exec.Command patterns instead of shell execution
  - Pattern matching for command arguments

### Runtime Security

- **Input Validation**: Request validation middleware to validate incoming requests
- **Rate Limiting**: Configurable rate limiting with two implementations:
  - Redis-based for distributed systems
  - In-memory for standalone deployments
- **Security Headers**: Middleware for security headers including:
  - Content-Security-Policy
  - X-XSS-Protection
  - X-Content-Type-Options
  - X-Frame-Options
  - Referrer-Policy
- **CORS Configuration**: Configurable CORS middleware with safe defaults
- **JWT Authentication**: Secure implementation of JWT with:
  - Token expiration
  - Refresh token rotation
  - Signature validation

### Dependency Security

- **Dependency Scanning**: Automated scanning of dependencies using `nancy` to detect:
  - Known vulnerabilities in dependencies
  - Outdated packages with security issues
  - Insecure dependency configurations

## Security Workflows

The security features are integrated into the development workflow through:

### Local Development

- **Pre-commit Hooks**: The project uses Husky to run pre-commit hooks:

  ```bash
  npm run security
  ```

  This runs `gosec` to identify security issues before code is committed.

- **Manual Security Scanning**: Developers can run security checks manually:

  ```bash
  # Full security scan
  npm run security

  # Focused scan on specific directories
  gosec ./internal/...

  # Dependency scanning
  go list -json -deps ./... | nancy sleuth
  ```

### CI/CD Integration

- **Pull Request Checks**: Security scans run automatically on pull requests:

  - Static code analysis with `gosec`
  - Dependency scanning with `nancy`
  - Results are uploaded to GitHub

- **Weekly Scheduled Scans**: Security scans run weekly to catch newly discovered vulnerabilities:
  ```yaml
  schedule:
    - cron: "0 0 * * 1" # Run weekly on Mondays
  ```

## Security Best Practices

The scaffolded applications follow these security best practices:

1. **Secure Configuration**: Sensitive configuration is loaded from environment variables
2. **Minimal Permissions**: Files and directories use minimal required permissions
3. **Error Handling**: Proper error handling without leaking sensitive information
4. **Secure Defaults**: Security features are enabled by default with secure configurations
5. **Dependency Management**: Regular updates of dependencies to address security issues
6. **Input Validation**: All user input is validated before processing
7. **Output Encoding**: Proper encoding of output to prevent XSS attacks
8. **Secure File Operations**: Safe file operations with atomic writes
9. **Logging**: Structured logging without sensitive information
10. **Authentication**: Strong authentication mechanisms with proper token management

## Adding Custom Security Features

You can extend the security features of the scaffolded application:

### Adding Custom Security Middleware

```go
// Example of adding a custom security middleware
func CustomSecurityMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Your security logic here
        c.Next()
    }
}

// In your router setup
router.Use(CustomSecurityMiddleware())
```

### Customizing Rate Limiting

```go
// In your configuration
type RateLimitConfig struct {
    RequestsPerMinute int
    BurstSize         int
    TTL               time.Duration
    KeyFunc           func(*gin.Context) string
}

// Create a custom key function based on IP and user ID
config.KeyFunc = func(c *gin.Context) string {
    userID := getUserIDFromContext(c)
    if userID != "" {
        return fmt.Sprintf("%s:%s", c.ClientIP(), userID)
    }
    return c.ClientIP()
}
```

## Security Reports

Security issues can be reported through:

1. GitHub Issues: For non-sensitive security issues
2. Email: For sensitive security issues, please email [security@example.com](mailto:security@example.com)

## Security Updates

The project will:

1. Regularly update dependencies to address security vulnerabilities
2. Release security patches for critical issues
3. Document security-related changes in release notes

## Further Reading

- [OWASP Go Security Cheatsheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Cheat_Sheet.html)
- [Go Security Best Practices](https://docs.google.com/document/d/1v9Hg9YLCVC4idrG4uHs7tV-R-QS9u6wKbX4Dj1lWGeg/edit)
- [Gosec Documentation](https://github.com/securego/gosec)
