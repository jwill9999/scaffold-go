package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateGenerator struct {
	ProjectName string
	Module      string
	Features    map[string]bool
	Config      interface{}
	Resources   []Resource
	Templates   map[string]string
}

type Resource struct {
	Name     string
	Resource string
}

func NewTemplateGenerator(name, module string, features []string, config interface{}) *TemplateGenerator {
	featureMap := make(map[string]bool)
	for _, f := range features {
		featureMap[f] = true
	}

	return &TemplateGenerator{
		ProjectName: name,
		Module:      module,
		Features:    featureMap,
		Config:      config,
		Resources:   []Resource{},
		Templates: map[string]string{
			// Core application files
			"cmd/api/main.go":           "tools/scaffold/templates/main.go.tmpl",
			"internal/config/config.go": "tools/scaffold/templates/config.go.tmpl",
			"Dockerfile":                "tools/scaffold/templates/Dockerfile.tmpl",
			"docker-compose.yml":        "tools/scaffold/templates/docker-compose.yml.tmpl",

			// Essential packages for a basic application
			"pkg/logger/logger.go":              "tools/scaffold/templates/logger.go.tmpl",
			"pkg/database/database.go":          "tools/scaffold/templates/database.go.tmpl",
			"internal/handlers/handlers.go":     "tools/scaffold/templates/handlers.go.tmpl",
			"internal/repository/repository.go": "tools/scaffold/templates/repository.go.tmpl",
			"internal/services/services.go":     "tools/scaffold/templates/service.go.tmpl",
			"internal/models/models.go":         "tools/scaffold/templates/model.go.tmpl",
		},
	}
}

func (g *TemplateGenerator) Generate() error {
	for filename, templatePath := range g.Templates {
		if err := g.generateFile(filename, templatePath); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}
	return nil
}

func (g *TemplateGenerator) generateFile(filename, templatePath string) error {
	// Read template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Create output file
	outputPath := filepath.Join(g.ProjectName, filename)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			// If we're already returning with an error, don't overwrite it
			if err == nil {
				err = fmt.Errorf("failed to close file %s: %w", outputPath, closeErr)
			}
		}
	}()

	// Execute template
	data := struct {
		ProjectName string
		Name        string
		Module      string
		Features    map[string]bool
		Config      interface{}
		Resources   []Resource
	}{
		ProjectName: g.ProjectName,
		Name:        g.ProjectName,
		Module:      g.Module,
		Features:    g.Features,
		Config:      g.Config,
		Resources:   g.Resources,
	}

	if err := tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}

func (g *TemplateGenerator) AddTemplate(name, path string) {
	g.Templates[name] = path
}

func (g *TemplateGenerator) RemoveTemplate(name string) {
	delete(g.Templates, name)
}
