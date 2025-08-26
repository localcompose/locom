package stage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func Init(targetDir string) error {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("creating target folder: %w", err)
		}
	}

	empty, err := isDirEmpty(targetDir)
	if err != nil {
		return fmt.Errorf("checking directory: %w", err)
	}
	if !empty {
		return fmt.Errorf("directory %q is not empty; please run 'locom init' in an empty folder", targetDir)
	}

	locomDir := filepath.Join(targetDir, ".locom")
	if _, err := os.Stat(locomDir); err == nil {
		return fmt.Errorf("stage already initialized: %q already exists", locomDir)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking .locom directory: %w", err)
	}

	if err := os.Mkdir(locomDir, 0755); err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("creating .locom directory: %w", err)
	}

	ymlPath := filepath.Join(locomDir, "locom.yml")
	f, err := os.OpenFile(ymlPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating locom.yml: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(defaultConfig); err != nil {
		return fmt.Errorf("writing locom.yml: %w", err)
	}

	fmt.Printf("Initialized empty Locom stage in %s\n", filepath.Join(targetDir, ".locom/"))
	return nil
}

func isDirEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.Name()[0] == '.' {
			continue
		}
		return false, nil
	}
	return true, nil
}

const defaultConfig = `stage:
  network:
    name: locom

    bind:
      address: 127.0.0.1

    dns:
      suffix: .locom.self

    proxy:
      name: traefik
      type: 
        engine: traefik
        version: 2.10

apps:
`
