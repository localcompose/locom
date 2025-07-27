package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Stage struct {
		Network struct {
			Name string `yaml:"name"`
			Bind struct {
				Address string `yaml:"address"`
			} `yaml:"bind"`
			DNS struct {
				Suffix string `yaml:"suffix"`
			} `yaml:"dns"`
			Proxy struct {
				Name string `yaml:"name"`
				Type struct {
					Engine  string `yaml:"engine"`
					Version string `yaml:"version"`
				} `yaml:"type"`
			} `yaml:"proxy"`
		} `yaml:"network"`
	} `yaml:"stage"`

	Apps map[string]any `yaml:"apps"`
}

func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling yaml: %w", err)
	}
	return &cfg, nil
}
