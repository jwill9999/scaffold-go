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
	Templates   map[string]string
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
		Templates: map[string]string{
			"main.go":            "templates/main.go.tmpl",
			"config.go":          "templates/config.go.tmpl",
			"Dockerfile":         "templates/Dockerfile.tmpl",
			"docker-compose.yml": "templates/docker-compose.yml.tmpl",
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
	defer out.Close()

	// Execute template
	data := struct {
		ProjectName string
		Module      string
		Features    map[string]bool
		Config      interface{}
	}{
		ProjectName: g.ProjectName,
		Module:      g.Module,
		Features:    g.Features,
		Config:      g.Config,
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
