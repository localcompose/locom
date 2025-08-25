package compose_test

import (
	"testing"

	"github.com/localcompose/locom/internal/compose"
)

func TestGetTraefikCompose(t *testing.T) {
	const networkName = "locom-net"
	cfg := compose.GetTraefikCompose(networkName)

	if len(cfg.Services) == 0 {
		t.Fatal("expected at least one service defined")
	}

	traefikService, ok := cfg.Services["traefik"]
	if !ok {
		t.Fatal("expected 'traefik' service to be defined")
	}

	if traefikService.ContainerName != "traefik" {
		t.Errorf("expected container name 'traefik', got %q", traefikService.ContainerName)
	}

	if len(cfg.Networks) == 0 {
		t.Fatal("expected at least one network defined")
	}

	if net, ok := cfg.Networks[networkName]; !ok {
		t.Errorf("expected network %q to be defined", networkName)
	} else if !net.External {
		t.Errorf("expected network %q to be external", networkName)
	}
}
