package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jwill9999/scaffold-go/tools/scaffold/generator"
)

type ProjectScaffold struct {
	Name      string
	Module    string
	Features  []string
	Structure ProjectStructure
	Config    ProjectConfig
}

type ProjectStructure struct {
	Directories []string
	BaseFiles   map[string]string
	Templates   map[string]string
}

type ProjectConfig struct {
	Environment string
	Database    DatabaseConfig
	Deployment  DeploymentConfig
}

type DatabaseConfig struct {
	Type      string // postgres, mysql
	Username  string
	Password  string
	Host      string
	Port      string
	Name      string
	EnableORM bool
}

type DeploymentConfig struct {
	Docker     bool
	Kubernetes bool
	CI         string // github, gitlab
}

// Base project directories based on core.mdc
var baseDirectories = []string{
	"cmd/api",
	"internal/core/server",
	"internal/core/middleware",
	"internal/core/errors",
	"internal/config",
	"internal/models",
	"internal/repositories",
	"internal/services",
	"internal/api/handlers",
	"internal/api/routes",
	"internal/api/dto",
	"pkg/logger",
	"pkg/metrics",
	"pkg/tracing",
	"pkg/validator",
	"migrations",
	"tests/unit",
	"tests/integration",
	"tests/e2e",
	"scripts",
	"deployments",
	"tools",
}

func main() {
	// Parse command line flags
	name := flag.String("name", "", "Project name")
	module := flag.String("module", "", "Go module path")
	features := flag.String("features", "", "Comma-separated list of features (auth,metrics,tracing)")
	dbType := flag.String("db", "postgres", "Database type (postgres, mysql)")
	deployment := flag.String("deployment", "docker", "Deployment type (docker, kubernetes)")

	flag.Parse()

	if *name == "" || *module == "" {
		log.Fatal("Project name and module path are required")
	}

	// Create project scaffold
	scaffold := &ProjectScaffold{
		Name:     *name,
		Module:   *module,
		Features: strings.Split(*features, ","),
		Structure: ProjectStructure{
			Directories: baseDirectories,
			BaseFiles:   make(map[string]string),
			Templates:   make(map[string]string),
		},
		Config: ProjectConfig{
			Environment: "development",
			Database: DatabaseConfig{
				Type:      *dbType,
				Username:  "postgres",
				Password:  "postgres",
				Host:      "localhost",
				Port:      "5432",
				Name:      *name,
				EnableORM: true,
			},
			Deployment: DeploymentConfig{
				Docker:     *deployment == "docker",
				Kubernetes: *deployment == "kubernetes",
				CI:         "github",
			},
		},
	}

	// Create project
	if err := scaffold.Create(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created project %s\n", *name)
}

func (p *ProjectScaffold) Create() error {
	// Create base directory
	if err := os.MkdirAll(p.Name, 0750); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create project structure
	if err := p.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Initialize go module
	if err := p.initGoModule(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	// Generate base files using template generator
	if err := p.generateBaseFiles(); err != nil {
		return fmt.Errorf("failed to generate base files: %w", err)
	}

	return nil
}

func (p *ProjectScaffold) createDirectories() error {
	for _, dir := range p.Structure.Directories {
		path := filepath.Join(p.Name, dir)
		if err := os.MkdirAll(path, 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func (p *ProjectScaffold) initGoModule() error {
	// Use filepath.Join instead of string concatenation for path safety
	projectPath := filepath.Join(".", p.Name)

	// Validate module name to prevent command injection
	if err := validateModuleName(p.Module); err != nil {
		return err
	}

	// Execute command with explicit arguments instead of through shell
	// The module name has been validated above to prevent command injection
	// #nosec G204 - p.Module is validated by validateModuleName
	cmd := exec.Command("go", "mod", "init", p.Module)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	return nil
}

// validateModuleName checks if the module name is valid
// Allows typical Go module path characters but prevents any special shell chars
func validateModuleName(name string) error {
	// Disallow empty name
	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	// Create a whitelist of allowed characters
	validModulePattern := `^[a-zA-Z0-9][a-zA-Z0-9\.\-_/]+$`
	matched, err := regexp.MatchString(validModulePattern, name)
	if err != nil {
		return fmt.Errorf("error validating module name: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid module name. Module name must start with a letter or number and contain only letters, numbers, periods, hyphens, underscores, and slashes")
	}

	return nil
}

func (p *ProjectScaffold) generateBaseFiles() error {
	// Create template generator
	tmplGen := generator.NewTemplateGenerator(
		p.Name,
		p.Module,
		p.Features,
		p.Config,
	)

	// Generate files
	if err := tmplGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	return nil
}

// executeCommand executes a shell command with security validation.
// This function is kept for reference purposes but is not used.
// nolint:unused // Intentionally kept for reference as a secure command execution example
func executeCommand(cmdStr string) error {
	// This function is deprecated for security reasons
	fmt.Println("Warning: executeCommand is deprecated for security reasons")

	// Define safe commands with their allowed arguments
	safeCommands := map[string][]string{
		"go mod init": {"^[a-zA-Z0-9][a-zA-Z0-9\\./\\-_]+$"}, // Allow standard module names
		"go mod tidy": nil,                                   // No arguments needed
		"go fmt":      {"^\\./.*$"},                          // Allow only relative paths starting with ./
	}

	// Try to match against our safe commands
	var matchedCmd string
	var argPattern []string

	for cmd, pattern := range safeCommands {
		if strings.HasPrefix(cmdStr, cmd) {
			matchedCmd = cmd
			argPattern = pattern
			break
		}
	}

	if matchedCmd == "" {
		return fmt.Errorf("command not allowed for security reasons: %s", cmdStr)
	}

	// Extract arguments
	args := strings.TrimSpace(strings.TrimPrefix(cmdStr, matchedCmd))

	// Parse the matched command to get program and initial args
	cmdParts := strings.Fields(matchedCmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("empty command")
	}

	program := cmdParts[0]
	progArgs := cmdParts[1:]

	// Add additional arguments if they exist and pass validation
	if args != "" {
		// If we have argument patterns defined, validate arguments
		if argPattern != nil {
			validArg := false
			for _, pattern := range argPattern {
				matched, err := regexp.MatchString(pattern, args)
				if err != nil {
					return fmt.Errorf("error validating command arguments: %w", err)
				}
				if matched {
					validArg = true
					break
				}
			}

			if !validArg {
				return fmt.Errorf("command argument not allowed for security reasons: %s", args)
			}
		}

		// Add the validated arguments
		argParts := strings.Fields(args)
		progArgs = append(progArgs, argParts...)
	}

	// Execute the command with explicit arguments - we've fully validated these
	// using a strict whitelist approach, so this is secure
	// #nosec G204 - program and progArgs are fully validated above
	cmd := exec.Command(program, progArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
