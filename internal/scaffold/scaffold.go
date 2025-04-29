package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jwill9999/scaffold-go/pkg/logger"
)

//go:embed templates/api/* templates/config/* templates/docker/*
var templateFS embed.FS

type Generator struct {
	ProjectName string
	ModulePath  string
	Features    []string
	DBType      string
	Deployment  string
	OutputDir   string
	Logger      *logger.Logger
}

func NewGenerator(name, module, features, dbType, deployment, outputDir string, logger *logger.Logger) *Generator {
	return &Generator{
		ProjectName: name,
		ModulePath:  module,
		Features:    parseFeatures(features),
		DBType:      dbType,
		Deployment:  deployment,
		OutputDir:   outputDir,
		Logger:      logger,
	}
}

func parseFeatures(features string) []string {
	if features == "" {
		return []string{}
	}
	return strings.Split(features, ",")
}

func (g *Generator) Generate() error {
	// Create project directory
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate project structure
	if err := g.generateProjectStructure(); err != nil {
		return fmt.Errorf("failed to generate project structure: %w", err)
	}

	// Generate base files
	if err := g.generateBaseFiles(); err != nil {
		return fmt.Errorf("failed to generate base files: %w", err)
	}

	// Generate feature-specific files
	if err := g.generateFeatures(); err != nil {
		return fmt.Errorf("failed to generate features: %w", err)
	}

	return nil
}

func (g *Generator) generateProjectStructure() error {
	dirs := []string{
		"cmd/api",
		"internal/config",
		"internal/handlers",
		"internal/models",
		"internal/repository",
		"internal/services",
		"internal/core/middleware",
		"internal/core/errors",
		"internal/core/server",
		"pkg/database",
		"pkg/logger",
		"pkg/security",
		"pkg/metrics",
	}

	for _, dir := range dirs {
		path := filepath.Join(g.OutputDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (g *Generator) generateBaseFiles() error {
	files := map[string]string{
		"main.go":            "templates/api/main.go.tmpl",
		"go.mod":             "templates/config/go.mod.tmpl",
		"config.yaml":        "templates/config/config.yaml.tmpl",
		"Dockerfile":         "templates/docker/Dockerfile.tmpl",
		"docker-compose.yml": "templates/docker/docker-compose.tmpl",
	}

	for target, tmpl := range files {
		if err := g.generateFile(target, tmpl); err != nil {
			return fmt.Errorf("failed to generate %s: %w", target, err)
		}
	}

	return nil
}

func (g *Generator) generateFeatures() error {
	for _, feature := range g.Features {
		if err := g.generateFeature(feature); err != nil {
			return fmt.Errorf("failed to generate feature %s: %w", feature, err)
		}
	}
	return nil
}

func (g *Generator) generateFile(target, tmpl string) error {
	content, err := templateFS.ReadFile(tmpl)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmpl, err)
	}

	t, err := template.New(filepath.Base(tmpl)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmpl, err)
	}

	targetPath := filepath.Join(g.OutputDir, target)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", target, err)
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", target, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", target, cerr)
		}
	}()

	if err := t.Execute(f, g); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", tmpl, err)
	}

	return err
}

func (g *Generator) generateFeature(feature string) error {
	switch feature {
	case "auth":
		return g.generateAuthFeature()
	case "metrics":
		return g.generateMetricsFeature()
	case "tracing":
		return g.generateTracingFeature()
	default:
		return fmt.Errorf("unknown feature: %s", feature)
	}
}

// Feature generation methods would go here
func (g *Generator) generateAuthFeature() error {
	// TODO: Implement auth feature generation
	return nil
}

func (g *Generator) generateMetricsFeature() error {
	// TODO: Implement metrics feature generation
	return nil
}

func (g *Generator) generateTracingFeature() error {
	// TODO: Implement tracing feature generation
	return nil
}
