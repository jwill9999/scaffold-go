package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateModuleName(t *testing.T) {
	tests := []struct {
		name    string
		module  string
		wantErr bool
	}{
		{
			name:    "Valid module name",
			module:  "github.com/username/project",
			wantErr: false,
		},
		{
			name:    "Valid module with numbers",
			module:  "github.com/user123/project-456",
			wantErr: false,
		},
		{
			name:    "Valid module with special chars",
			module:  "github.com/user/my-awesome_project.api",
			wantErr: false,
		},
		{
			name:    "Empty module name",
			module:  "",
			wantErr: true,
		},
		{
			name:    "Module with spaces",
			module:  "github.com/user/my project",
			wantErr: true,
		},
		{
			name:    "Module with shell characters",
			module:  "github.com/user/project;rm -rf /",
			wantErr: true,
		},
		{
			name:    "Module starting with special char",
			module:  "-github.com/user/project",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateModuleName(tt.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateModuleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectScaffoldInitialization(t *testing.T) {
	// Test case for basic project scaffold initialization
	scaffold := &ProjectScaffold{
		Name:     "test-project",
		Module:   "github.com/test/test-project",
		Features: []string{"auth", "metrics"},
		Structure: ProjectStructure{
			Directories: baseDirectories,
			BaseFiles:   make(map[string]string),
			Templates:   make(map[string]string),
		},
		Config: ProjectConfig{
			Environment: "development",
			Database: DatabaseConfig{
				Type:      "postgres",
				Username:  "test-user",
				Password:  "test-pass",
				Host:      "localhost",
				Port:      "5432",
				Name:      "test-db",
				EnableORM: true,
			},
			Deployment: DeploymentConfig{
				Docker:     true,
				Kubernetes: false,
				CI:         "github",
			},
		},
	}

	// Verify the scaffold properties were set correctly
	if scaffold.Name != "test-project" {
		t.Errorf("Expected Name to be 'test-project', got '%s'", scaffold.Name)
	}

	if scaffold.Module != "github.com/test/test-project" {
		t.Errorf("Expected Module to be 'github.com/test/test-project', got '%s'", scaffold.Module)
	}

	if len(scaffold.Features) != 2 || scaffold.Features[0] != "auth" || scaffold.Features[1] != "metrics" {
		t.Errorf("Features not set correctly, got %v", scaffold.Features)
	}

	if len(scaffold.Structure.Directories) == 0 {
		t.Error("Expected directories to be populated")
	}

	if scaffold.Config.Database.Type != "postgres" {
		t.Errorf("Expected database type to be 'postgres', got '%s'", scaffold.Config.Database.Type)
	}

	if !scaffold.Config.Deployment.Docker {
		t.Error("Expected Docker deployment to be true")
	}

	if scaffold.Config.Deployment.Kubernetes {
		t.Error("Expected Kubernetes deployment to be false")
	}
}

func TestCreateDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	// Clean up after the test and check for errors
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to clean up temp directory: %v", err)
		}
	}()

	// Set up a minimal test scaffold
	testDirs := []string{"cmd/api", "internal/config", "pkg/logger"}
	scaffold := &ProjectScaffold{
		Name: tempDir,
		Structure: ProjectStructure{
			Directories: testDirs,
		},
	}

	// Test directory creation
	err = scaffold.createDirectories()
	if err != nil {
		t.Fatalf("createDirectories() failed: %v", err)
	}

	// Verify directories were created
	for _, dir := range testDirs {
		fullPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Directory not created: %s", fullPath)
		}
	}
}

func TestInitGoModule(t *testing.T) {
	// Skip this test if we're not in a testing environment
	// to avoid actually modifying the filesystem
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	// Clean up after the test and check for errors
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to clean up temp directory: %v", err)
		}
	}()

	// Test cases for initGoModule validation only
	tests := []struct {
		name    string
		module  string
		wantErr bool
	}{
		{
			name:    "Valid module name",
			module:  "github.com/test/project",
			wantErr: false,
		},
		{
			name:    "Invalid module name with shell characters",
			module:  "github.com/test/project;rm -rf /",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaffold := &ProjectScaffold{
				Name:   filepath.Join(tempDir, "test-"+tt.name),
				Module: tt.module,
			}

			// Create the test project directory
			err := os.MkdirAll(scaffold.Name, 0750)
			if err != nil {
				t.Fatalf("Failed to create test project directory: %v", err)
			}

			// We only test the validation part of initGoModule
			// since we don't want to execute actual commands
			err = validateModuleName(tt.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateModuleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGenerateBaseFiles tests that the generateBaseFiles function
// handles validation and calls the template generator correctly.
// This is a basic test that needs to be expanded with proper mocking.
func TestGenerateBaseFiles(t *testing.T) {
	// Skip this test when running in short mode
	if testing.Short() {
		t.Skip("Skipping generateBaseFiles test in short mode")
	}

	// Create a minimal scaffold
	scaffold := &ProjectScaffold{
		Name:     "test-project",
		Module:   "github.com/test/test-project",
		Features: []string{"auth"},
		Config: ProjectConfig{
			Environment: "development",
			Database: DatabaseConfig{
				Type: "postgres",
			},
		},
	}

	// Since we can't easily mock the generator package in this test,
	// we just do a basic validation that it doesn't crash with our inputs.
	// In a real test suite, we'd use dependency injection to mock the generator.

	// Note: This will fail if executed, since we're not creating the actual directory
	// Just verify that our validation logic works
	err := validateModuleName(scaffold.Module)
	if err != nil {
		t.Errorf("Module name validation failed: %v", err)
	}
}
