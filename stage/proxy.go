package stage

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/localcompose/locom/compose"
	"github.com/localcompose/locom/config"
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

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating proxy folder: %w", err)
	}

	composeData := compose.GetTraefikCompose(networkName)

	ymlData, err := yaml.Marshal(composeData)
	if err != nil {
		return fmt.Errorf("serializing yaml: %w", err)
	}

	filePath := filepath.Join(targetDir, "docker-compose.yml")
	if err := os.WriteFile(filePath, ymlData, 0644); err != nil {
		return fmt.Errorf("writing docker-compose.yml: %w", err)
	}

	return nil
}
