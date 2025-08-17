package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/localcompose/locom/config"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	yamlData := `
stage:
  network:
    name: testnet
    bind:
      address: 127.0.0.1
    dns:
      suffix: .locom.self
    proxy:
      name: traefik
      type:
        engine: traefik
        version: 2.10
apps: {}
`

	tmpFile := filepath.Join(t.TempDir(), "locom.yml")
	require.NoError(t, os.WriteFile(tmpFile, []byte(yamlData), 0644))

	cfg, err := config.LoadConfig(tmpFile)
	require.NoError(t, err)
	require.Equal(t, "testnet", cfg.Stage.Network.Name)
}
