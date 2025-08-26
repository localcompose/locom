package stage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit_CreatesStageInEmptyDir(t *testing.T) {
	tmp := t.TempDir()

	err := Init(tmp)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	locomDir := filepath.Join(tmp, ".locom")
	ymlPath := filepath.Join(locomDir, "locom.yml")

	if _, err := os.Stat(locomDir); err != nil {
		t.Errorf(".locom directory was not created: %v", err)
	}
	if _, err := os.Stat(ymlPath); err != nil {
		t.Errorf("locom.yml was not created: %v", err)
	}
}

func TestInit_FailsInNonEmptyDir(t *testing.T) {
	tmp := t.TempDir()

	// Add a dummy file
	if err := os.WriteFile(filepath.Join(tmp, "README.md"), []byte("test"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err := Init(tmp)
	if err == nil {
		t.Fatal("expected Init to fail in non-empty directory, but it succeeded")
	}
}

func TestInit_FailsIfAlreadyInitialized(t *testing.T) {
	tmp := t.TempDir()

	// First init
	if err := Init(tmp); err != nil {
		t.Fatalf("first Init failed: %v", err)
	}

	// Second init should fail
	err := Init(tmp)
	if err == nil {
		t.Fatal("expected Init to fail if .locom already exists, but it succeeded")
	}
}
