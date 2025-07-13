package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdNew)
}

var cmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a new configuration",
	Run: runNew,
}

func runNew(cmd *cobra.Command, args []string) {
	path := "locom.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("Failed to create config: %v\n", err)
			return
		}
		defer file.Close()
		file.WriteString("# locom configuration\nversion: 0.0.1\n")
		fmt.Printf("Created new config at %s\n", filepath.Join(".", path))
	} else {
		fmt.Printf("Config already exists at %s\n", filepath.Join(".", path))
	}
}
