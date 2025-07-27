package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(cmdProxy)
}

var cmdProxy = &cobra.Command{
	Use:   "proxy",
	Short: "Create a default docker-compose configuration with Traefik proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "proxy"
		return runProxy(target)
	},
}

func runProxy(targetDir string) error {
	configPath := filepath.Join(".locom", "locom.yml")
	cfg, err := loadConfig(configPath)
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

	compose := getTraefikComposeWithNetwork(networkName)

	ymlData, err := yaml.Marshal(compose)
	if err != nil {
		return fmt.Errorf("serializing yaml: %w", err)
	}

	filePath := filepath.Join(targetDir, "docker-compose.yml")
	if err := os.WriteFile(filePath, ymlData, 0644); err != nil {
		return fmt.Errorf("writing docker-compose.yml: %w", err)
	}

	fmt.Printf("Docker Compose configuration written to %s\n", filePath)
	return nil
}

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}

	return &cfg, nil
}

type Config struct {
	Stage struct {
		Network struct {
			Name string `yaml:"name"`
		} `yaml:"network"`
	} `yaml:"stage"`
}
func getTraefikComposeWithNetwork(network string) map[string]interface{} {
	return map[string]interface{}{
		"networks": map[string]interface{}{
			network: map[string]interface{}{
				"external": true,
			},
		},
		"services": map[string]interface{}{
			"traefik": map[string]interface{}{
				"image":         "traefik:v2.10",
				"container_name": "traefik",
				"restart":       "unless-stopped",
				"command": []string{
					"--api.dashboard=true",
					"--api.insecure=true",
					"--providers.docker=true",
					"--providers.docker.exposedbydefault=false",
					"--entrypoints.web.address=:80",
					"--entrypoints.websecure.address=:443",
					"--providers.file.directory=/etc/traefik/dynamic",
					"--providers.file.watch=true",
				},
				"ports": []string{
					"80:80",
					"443:443",
					"8080:8080",
				},
				"volumes": []string{
					"/var/run/docker.sock:/var/run/docker.sock:ro",
					"./config:/etc/traefik/dynamic",
				},
				"networks": []string{
					network,
				},
				"labels": map[string]string{
					"traefik.enable":                           "true",
					"traefik.http.routers.traefik.rule":        "Host(`proxy.locom.self`)",
					"traefik.http.routers.traefik.service":     "api@internal",
					"traefik.http.routers.traefik.entrypoints": "web",
				},
			},
		},
	}
}
