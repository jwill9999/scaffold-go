# Testing Strategy for Go Scaffolding System

This document outlines our comprehensive testing approach for the Go API scaffolding system to ensure reliability, correctness, and maintainability.

## Testing Levels

### 1. Unit Tests

Unit tests focus on validating individual components of the scaffolding system:

- **Template parsing and rendering**
- **Configuration loading and validation**
- **File path handling and generation**
- **Variable substitution**

**Example:**
```go
// tools/scaffold/generator/parser_test.go
func TestTemplateParser(t *testing.T) {
    parser := NewParser()
    result := parser.Parse("{{ .ProjectName }}")
    assert.NotNil(t, result)
}
```

### 2. Template Validation Tests

Ensures all templates are syntactically correct and produce the expected output:

- **Syntax checking** for all templates
- **Golden file testing** to compare generated output with expected results
- **Variable interpolation** tests for different input configurations

**Example:**
```go
// tools/scaffold/templates/validator_test.go
func TestTemplatesAreValid(t *testing.T) {
    templates, err := filepath.Glob("../../templates/*.tmpl")
    assert.NoError(t, err)
    
    for _, tmpl := range templates {
        _, err := template.ParseFiles(tmpl)
        assert.NoError(t, err, "Template %s should be valid", tmpl)
    }
}
```

### 3. Integration Tests

Tests the scaffolding system as a whole:

- **End-to-end project generation**
- **Command-line interface** testing
- **Plugin interactions**
- **Configuration combinations**

**Example:**
```go
func TestGenerateBasicProject(t *testing.T) {
    tmpDir, _ := ioutil.TempDir("", "scaffold-test-")
    defer os.RemoveAll(tmpDir)
    
    args := []string{"init", "--name", "testapi", "--module", "github.com/test/testapi"}
    err := RunScaffold(args, tmpDir)
    assert.NoError(t, err)
    
    // Check for expected files
    assert.FileExists(t, filepath.Join(tmpDir, "main.go"))
    assert.FileExists(t, filepath.Join(tmpDir, "go.mod"))
}
```

### 4. Compilation Tests

Verifies that generated projects compile and run correctly:

- **Build tests** for different project configurations
- **Import resolution** checks
- **Dependency management** validation

**Example:**
```go
func TestGeneratedProjectCompiles(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping compilation test in short mode")
    }
    
    tmpDir, _ := ioutil.TempDir("", "scaffold-test-")
    defer os.RemoveAll(tmpDir)
    
    RunScaffold([]string{"init", "--name", "testapi"}, tmpDir)
    
    cmd := exec.Command("go", "build")
    cmd.Dir = tmpDir
    output, err := cmd.CombinedOutput()
    assert.NoError(t, err, "Generated project should compile: %s", output)
}
```

### 5. Plugin Tests

Dedicated tests for each plugin:

- **Database plugins** (PostgreSQL, MySQL, etc.)
- **Authentication plugins** (JWT, OAuth2, etc.)
- **Caching plugins** (Redis, in-memory, etc.)
- **Feature plugins** (pagination, filtering, etc.)

**Example:**
```go
func TestDatabasePlugin(t *testing.T) {
    tmpDir := setupTestProject(t)
    defer os.RemoveAll(tmpDir)
    
    err := RunScaffold([]string{"plugin", "add", "--database", "postgres"}, tmpDir)
    assert.NoError(t, err)
    
    // Check for database-specific files
    assert.FileExists(t, filepath.Join(tmpDir, "pkg", "database", "postgres.go"))
}
```

## Test Fixtures and Tools

### Test Data

Standard test fixtures for consistent testing:

```go
// testdata/fixtures.go
var TestProjectData = &ProjectData{
    Name: "testapi",
    Module: "github.com/test/testapi",
    Database: "postgres",
    Auth: "jwt",
}
```

### Golden Files

Reference files for comparing generated output:

- Stored in `testdata/golden/`
- One file for each template/configuration combination
- Updated with `-update` flag when intentional changes are made

### CI Integration

- **Pre-commit hooks** run unit tests and linting
- **Pre-push hooks** run integration tests
- **GitHub Actions** run the full test suite, including compilation tests

## Test Organization

```
scaffold-go/
├── tools/
│   └── scaffold/
│       ├── generator/
│       │   └── *_test.go        # Unit tests
│       ├── testdata/            
│       │   ├── fixtures/        # Test fixtures
│       │   └── golden/          # Golden files
│       └── integration_test.go  # Integration tests
└── test/
    └── e2e/                     # End-to-end tests
```

## Test Reporting and Coverage

### Coverage Reports

Our testing infrastructure automatically generates coverage reports in multiple formats in the `test-output` directory:

- **HTML Coverage Report**: Generated with `npm run test:html`, creates a visual representation of code coverage in `test-output/coverage.html`
- **JSON Test Reports**: Generated with `npm run test:ci`, creates structured test data in `test-output/test-report.json`
- **Coverage Data**: Both commands generate raw coverage data in `test-output/cover.out`

### Automated Reporting

Coverage and test results are integrated into our development workflow:

1. **Pre-push Hook**: Generates coverage report before pushing code
2. **GitHub Actions**: Runs tests and uploads coverage artifacts
3. **Codecov Integration**: Publishes coverage data to Codecov for tracking

### Viewing Reports

- Coverage reports are available as artifacts in GitHub Actions runs
- The Codecov dashboard shows historical coverage trends
- Locally, run `npm run test:html` and open `test-output/coverage.html` in a browser

### Coverage Goals

- **Core Utilities**: Aim for >90% coverage
- **Templates**: 100% syntax validation
- **Overall Project**: Maintain >75% coverage

## Implementation Plan

1. **Phase 1: Foundation**
   - Set up test infrastructure
   - Implement unit tests for core components
   - Add template validation tests

2. **Phase 2: Integration**
   - Implement golden file testing
   - Add basic integration tests
   - Set up CI pipeline for tests

3. **Phase 3: Expansion**
   - Add compilation tests
   - Implement plugin-specific tests
   - Create end-to-end test suite

4. **Phase 4: Automation**
   - Automate test fixture generation
   - Add performance benchmarks
   - Implement test coverage reporting

## Testing Best Practices

1. **Write tests first** when adding new features
2. **Maintain test fixtures** as the system evolves
3. **Update golden files** only when changes are intentional
4. **Run full test suite** before releases
5. **Keep tests fast** for developer productivity
6. **Test edge cases** and error conditions 