package stage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/localcompose/locom/stage"
)

func TestGenerateProxyComposeFiles(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".locom")
	targetDir := filepath.Join(tmpDir, "proxy")

	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Write minimal config file
	configContent := `
stage:
  network:
    name: testnet
`
	configPath := filepath.Join(configDir, "locom.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Run the function
	err := stage.GenerateProxyComposeFiles(configPath, targetDir)
	if err != nil {
		t.Fatalf("GenerateProxyComposeFiles failed: %v", err)
	}

	sourcePath := filepath.Join(configDir, "proxy", "docker-compose.yml")
	targetPath := filepath.Join(tmpDir, "proxy", "docker-compose.yml")

	for _, path := range []string{sourcePath, targetPath} {
		// Check file was created
		// Basic sanity check on YAML content
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("missing file: %s: %v", path, err)
		} else if !containsAll(string(data), "traefik", "testnet") {

			t.Errorf("unexpected content in %s:\n%s", path, data)
		}
	}
}

func containsAll(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(s) > len(sub) && (s[0:len(sub)] == sub || contains(s[1:], sub))))
}
