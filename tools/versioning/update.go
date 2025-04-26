// tools/versioning/update.go
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: update-version <new_version>")
		os.Exit(1)
	}

	version := os.Args[1]

	// Validate version format with regex (e.g., v1.0.0)
	validVersion := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?$`)
	if !validVersion.MatchString(version) {
		fmt.Printf("Error: Invalid version format '%s'. Expected format: v1.0.0 or 1.0.0\n", version)
		os.Exit(1)
	}

	versionPath := "tools/scaffold/templates/VERSION"

	// Check if file exists before attempting to read it
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		fmt.Printf("Error: VERSION file not found at %s\n", versionPath)
		os.Exit(1)
	}

	currentContent, err := os.ReadFile(versionPath)
	if err != nil {
		fmt.Printf("Error reading VERSION file: %v\n", err)
		os.Exit(1)
	}

	currentStr := string(currentContent)
	today := time.Now().Format("2006-01-02")

	var releaseDate string
	if strings.Contains(currentStr, "Release Date:") {
		parts := strings.Split(currentStr, "Release Date:")
		dateLineParts := strings.Split(parts[1], "\n")
		releaseDate = strings.TrimSpace(dateLineParts[0])
	} else {
		releaseDate = today
	}

	// Extract template versions section
	var templateSection string
	if strings.Contains(currentStr, "Template Versions:") {
		parts := strings.Split(currentStr, "Template Versions:\n")
		if len(parts) > 1 {
			lastUpdateParts := strings.Split(parts[1], "Last Updated:")
			templateSection = lastUpdateParts[0]
		}
	}

	newContent := fmt.Sprintf("version: %s\n\nTemplate Versions:\n%s\nLast Updated: %s\nRelease Date: %s",
		version,
		templateSection,
		today,
		releaseDate)

	// Create a temporary file first and then rename it to avoid partial writes
	tempFile := versionPath + ".tmp"
	err = os.WriteFile(tempFile, []byte(newContent), 0600)
	if err != nil {
		fmt.Printf("Error writing temporary VERSION file: %v\n", err)
		os.Exit(1)
	}

	// Rename the temp file to the target file (atomic operation)
	if err := os.Rename(tempFile, versionPath); err != nil {
		fmt.Printf("Error updating VERSION file: %v\n", err)
		// Try to clean up the temp file
		if rmErr := os.Remove(tempFile); rmErr != nil {
			fmt.Printf("Additionally, failed to remove temporary file: %v\n", rmErr)
		}
		os.Exit(1)
	}

	fmt.Printf("Updated version to %s\n", version)
}
