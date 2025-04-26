package generator

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateGenerator struct {
	ProjectName string
	Module      string
	Features    map[string]bool
	Config      interface{}
	Resources   []Resource
	Templates   map[string]string
	// Base directory to prevent path traversal
	BaseDir string
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

	// Use current directory as the base directory
	baseDir, err := filepath.Abs(".")
	if err != nil {
		baseDir = "."
	}

	return &TemplateGenerator{
		ProjectName: name,
		Module:      module,
		Features:    featureMap,
		Config:      config,
		Resources:   []Resource{},
		BaseDir:     baseDir,
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

// sanitizePath prevents path traversal attacks by checking for ../ patterns and absolute paths
func (g *TemplateGenerator) sanitizePath(path string) (string, error) {
	// Check for .. patterns that might indicate path traversal
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("path contains illegal pattern: %s", path)
	}

	// If path is absolute, check that it's within base directory
	if filepath.IsAbs(path) {
		cleanPath := filepath.Clean(path)
		baseDir, err := filepath.Abs(g.BaseDir)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute base path: %w", err)
		}

		relPath, err := filepath.Rel(baseDir, cleanPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			return "", fmt.Errorf("path is outside base directory: %s", path)
		}

		return cleanPath, nil
	}

	// For relative paths, clean and return
	return filepath.Clean(path), nil
}

func (g *TemplateGenerator) generateFile(filename, templatePath string) error {
	// Sanitize paths
	sanitizedTemplatePath, err := g.sanitizePath(templatePath)
	if err != nil {
		return fmt.Errorf("invalid template path: %w", err)
	}

	// Validate that template exists
	if _, err := os.Stat(sanitizedTemplatePath); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", sanitizedTemplatePath)
	}

	// Read template
	tmpl, err := template.ParseFiles(sanitizedTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", sanitizedTemplatePath, err)
	}

	// Create output file
	// Ensure we create paths relative to project name to avoid path traversal
	if strings.Contains(filename, "..") || filepath.IsAbs(filename) {
		return fmt.Errorf("invalid output filename: %s", filename)
	}

	outputPath := filepath.Join(g.ProjectName, filename)
	outputDir := filepath.Dir(outputPath)

	// Create directories with secure permissions
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
	}

	// Use a more secure temporary file in the same directory
	// Use a cryptographically secure random suffix rather than a predictable one
	randomSuffix, err := generateSecureRandomString(8)
	if err != nil {
		return fmt.Errorf("failed to generate secure random string: %w", err)
	}

	tempFile := filepath.Join(outputDir, fmt.Sprintf(".tmp_%s_%s", filepath.Base(outputPath), randomSuffix))

	// Create file with secure permissions
	// Using os.OpenFile directly with strict permissions
	// #nosec G304 - We've already sanitized and validated the path above
	tempOut, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("failed to create temporary file for %s: %w", outputPath, err)
	}

	// Ensure tempOut is closed properly
	tempOutClosed := false
	defer func() {
		if !tempOutClosed {
			if closeErr := tempOut.Close(); closeErr != nil && err == nil {
				err = fmt.Errorf("failed to close temporary file: %w", closeErr)
			}
			// Always try to remove the temp file in case of an error
			if err != nil {
				_ = os.Remove(tempFile)
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

	if err := tmpl.Execute(tempOut, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", sanitizedTemplatePath, err)
	}

	// Close the temp file before renaming
	if err := tempOut.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}
	tempOutClosed = true

	// Verify the file was successfully written by checking its existence and size
	tempFileStat, err := os.Stat(tempFile)
	if err != nil {
		return fmt.Errorf("failed to stat temporary file: %w", err)
	}

	if tempFileStat.Size() == 0 {
		// Delete empty files and report error
		_ = os.Remove(tempFile)
		return fmt.Errorf("template generated empty file, generation failed")
	}

	// Rename the temp file to the final output file (atomic operation)
	if err := os.Rename(tempFile, outputPath); err != nil {
		// Try to remove the temp file if rename fails
		if rmErr := os.Remove(tempFile); rmErr != nil {
			return fmt.Errorf("failed to rename temp file: %v (cleanup error: %v)", err, rmErr)
		}
		return fmt.Errorf("failed to finalize file %s: %w", outputPath, err)
	}

	return nil
}

// generateSecureRandomString creates a cryptographically secure random string
// to be used for temporary file names
func generateSecureRandomString(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice to hold the random bytes
	bytes := make([]byte, length)

	// Read random bytes from crypto/rand
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert random bytes to letters
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes), nil
}

func (g *TemplateGenerator) AddTemplate(name, path string) {
	// Sanitize the path before adding
	sanitizedPath, err := g.sanitizePath(path)
	if err != nil {
		fmt.Printf("Warning: skipping invalid template path: %v\n", err)
		return
	}

	g.Templates[name] = sanitizedPath
}

func (g *TemplateGenerator) RemoveTemplate(name string) {
	delete(g.Templates, name)
}
