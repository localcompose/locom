package stage

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/localcompose/locom/internal/compose"
	"github.com/localcompose/locom/internal/config"
)

func GenerateProxyComposeFiles(configPath, targetDir string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	networkName := cfg.Stage.Network.Name
	if networkName == "" {
		return fmt.Errorf("network name not found in configuration")
	}

	// Generate the compose content
	composeData := compose.GetTraefikCompose(networkName)
	ymlData, err := yaml.Marshal(composeData)
	if err != nil {
		return fmt.Errorf("serializing yaml: %w", err)
	}

	// 1. Write the template source under .locom/proxy/
	sourcePath := filepath.Join(filepath.Dir(configPath), "proxy")
	if err := os.MkdirAll(sourcePath, 0755); err != nil {
		return fmt.Errorf("creating .locom/proxy folder: %w", err)
	}
	sourceFile := filepath.Join(sourcePath, "docker-compose.yml")
	if err := os.WriteFile(sourceFile, ymlData, 0644); err != nil {
		return fmt.Errorf("writing source docker-compose.yml: %w", err)
	}

	// 2. Copy it to ./proxy/docker-compose.yml if it doesn't already exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating proxy target folder: %w", err)
	}
	targetFile := filepath.Join(targetDir, "docker-compose.yml")

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		// Only write if file does not exist
		if err := os.WriteFile(targetFile, ymlData, 0644); err != nil {
			return fmt.Errorf("writing proxy/docker-compose.yml: %w", err)
		}
		fmt.Printf("Created %s from template\n", targetFile)
	} else {
		fmt.Printf("Skipped writing %s (already exists)\n", targetFile)
	}

	return nil
}
