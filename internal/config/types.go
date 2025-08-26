package config

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
