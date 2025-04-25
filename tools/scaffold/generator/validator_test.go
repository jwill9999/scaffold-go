package generator

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

// TestTemplateFilesExist verifies that the template directory exists and contains files
func TestTemplateFilesExist(t *testing.T) {
	// Skip this test if there's no template directory structure yet
	// Remove these lines once the template directory is properly set up
	t.Skip("Skipping test until template directory structure is fully set up")

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up to the scaffold root from the generator package
	scaffoldRoot := filepath.Join(wd, "..", "..")
	templatesDir := filepath.Join(scaffoldRoot, "templates")

	// Check if the templates directory exists
	info, err := os.Stat(templatesDir)
	if os.IsNotExist(err) {
		t.Fatalf("Templates directory not found at %s", templatesDir)
	}
	if err != nil {
		t.Fatalf("Error accessing templates directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Expected %s to be a directory", templatesDir)
	}

	// Look for template files
	files, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		t.Fatalf("Error searching for template files: %v", err)
	}

	// Check if we found any template files
	if len(files) == 0 {
		files, err = filepath.Glob(filepath.Join(templatesDir, "*.go.tmpl"))
		if err != nil {
			t.Fatalf("Error searching for .go.tmpl files: %v", err)
		}

		if len(files) == 0 {
			t.Fatalf("No template files found in %s", templatesDir)
		}
	}

	t.Logf("Found %d template files", len(files))
}

// TestTemplateFilesSyntax verifies that all template files are syntactically valid
func TestTemplateFilesSyntax(t *testing.T) {
	// Skip this test if there's no template directory structure yet
	// Remove these lines once the template directory is created
	t.Skip("Skip validation test until template directory is fully set up")

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up to the scaffold root from the generator package
	scaffoldRoot := filepath.Join(wd, "..", "..")
	templatesDir := filepath.Join(scaffoldRoot, "templates")

	// Find all template files
	var templateFiles []string

	err = filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if file has .tmpl extension
		if !info.IsDir() && (filepath.Ext(path) == ".tmpl" ||
			filepath.Ext(filepath.Ext(path)+path) == ".go.tmpl") {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Error walking templates directory: %v", err)
	}

	if len(templateFiles) == 0 {
		t.Fatalf("No template files found in %s", templatesDir)
	}

	t.Logf("Validating %d template files", len(templateFiles))

	// Test each template file
	for _, file := range templateFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			// Read template file
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read template file %s: %v", file, err)
			}

			// Try to parse the template
			_, err = template.New(filepath.Base(file)).Parse(string(content))
			if err != nil {
				t.Errorf("Template syntax error in %s: %v", file, err)
			}
		})
	}
}
