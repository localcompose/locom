package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdInit)
}

var cmdInit = &cobra.Command{
	Use:   "init [folder]",
	Short: "Initialize a new locom stage in the specified folder",
	Long: `Creates a .locom directory and a default locom.yml config file 
inside the given folder. Fails if the folder is not empty.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) == 1 {
			target = args[0]
		}
		return runInit(target)
	},
}

func runInit(targetDir string) error {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		// Folder does not exist: create it
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
		return fmt.Errorf("project already initialized: %q already exists", locomDir)
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

	fmt.Printf("Initialized empty Locom project in %s\n", filepath.Join(targetDir, ".locom/"))
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

const defaultConfig = `platform:
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
