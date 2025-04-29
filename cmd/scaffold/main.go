package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/jwill9999/scaffold-go/internal/scaffold"
	"github.com/jwill9999/scaffold-go/pkg/logger"
)

var (
	Version = "0.0.1"
)

type Config struct {
	Name       string
	Module     string
	Features   string
	DBType     string
	Deployment string
	OutputDir  string
}

func main() {
	config := Config{}
	log := logger.New(true)

	// Parse command line flags
	flag.StringVar(&config.Name, "name", "", "Project name")
	flag.StringVar(&config.Module, "module", "", "Go module path")
	flag.StringVar(&config.Features, "features", "", "Comma-separated list of features (auth,metrics,tracing)")
	flag.StringVar(&config.DBType, "db", "postgres", "Database type (postgres)")
	flag.StringVar(&config.Deployment, "deployment", "docker", "Deployment type (docker,k8s)")
	flag.StringVar(&config.OutputDir, "output", "", "Output directory (default: current directory)")
	flag.Parse()

	// Validate required flags
	if config.Name == "" || config.Module == "" {
		log.Error("Project name and module path are required")
		flag.Usage()
		os.Exit(1)
	}

	// Set default output directory if not specified
	if config.OutputDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Error("Failed to get current directory: %v", err)
			os.Exit(1)
		}
		config.OutputDir = filepath.Join(currentDir, config.Name)
	}

	// Create generator
	generator := scaffold.NewGenerator(
		config.Name,
		config.Module,
		config.Features,
		config.DBType,
		config.Deployment,
		config.OutputDir,
		log,
	)

	// Generate project
	if err := generator.Generate(); err != nil {
		log.Error("Failed to generate project: %v", err)
		os.Exit(1)
	}

	log.Info("Successfully generated project at %s", config.OutputDir)
}
