// tools/versioning/update.go
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: update-version <new_version>")
		os.Exit(1)
	}

	version := os.Args[1]
	versionPath := "tools/scaffold/templates/VERSION"

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

	err = os.WriteFile(versionPath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing VERSION file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated version to %s\n", version)
}
