package compose

type ComposeFile struct {
	Version  string                      `yaml:"version,omitempty"`
	Networks map[string]ExternalNetwork  `yaml:"networks"`
	Services map[string]Service          `yaml:"services"`
}

type Service struct {
	Image         string            `yaml:"image,omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
	Restart       string            `yaml:"restart,omitempty"`
	Command       []string          `yaml:"command,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Networks      []string          `yaml:"networks,omitempty"`
	Labels        map[string]string `yaml:"labels,omitempty"`
}

type ExternalNetwork struct {
	External bool `yaml:"external"`
}
