package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	if err := os.MkdirAll(p.Name, 0755); err != nil {
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
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func (p *ProjectScaffold) initGoModule() error {
	cmd := exec.Command("go", "mod", "init", p.Module)
	cmd.Dir = p.Name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
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

func executeCommand(cmd string) error {
	command := exec.Command("sh", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
