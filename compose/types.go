package compose

import "gopkg.in/yaml.v3"

// ComposeFile represents a Docker Compose file
type ComposeFile struct {
	Version  string                     `yaml:"version,omitempty"`
	Networks map[string]ExternalNetwork `yaml:"networks"`
	Services map[string]Service         `yaml:"services"`
}

// Service represents a single service in Compose
type Service struct {
	Image         string   `yaml:"image,omitempty"`
	ContainerName string   `yaml:"container_name,omitempty"`
	Restart       string   `yaml:"restart,omitempty"`
	Command       []string `yaml:"command,omitempty"`
	Ports         []string `yaml:"ports,omitempty"`
	Volumes       []string `yaml:"volumes,omitempty"`
	Networks      []string `yaml:"networks,omitempty"`

	// Use either Labels (unordered map) or LabelsNode (ordered with optional anchors)
	// Labels     map[string]string `yaml:"labels,omitempty"`
	LabelsNode *yaml.Node `yaml:"labels,omitempty"`
}

// ExternalNetwork represents an external network
type ExternalNetwork struct {
	External bool `yaml:"external"`
}
