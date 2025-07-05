package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "locom",
	Short: "locom manages local Docker Compose stacks",
	Long:  `locom is a CLI tool for managing local Docker Compose stacks in a minimal, offline-friendly way.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
