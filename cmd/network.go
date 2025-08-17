package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/localcompose/locom/config"
)

var cmdNetwork = &cobra.Command{
	Use:   "network",
	Short: "Ensure the Docker network defined in .locom/locom.yml exists",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := filepath.Join(".locom", "locom.yml")
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		networkName := cfg.Stage.Network.Name
		if networkName == "" {
			return fmt.Errorf("no network name found in config")
		}

		return ensureDockerNetwork(networkName)
	},
}

func init() {
	rootCmd.AddCommand(cmdNetwork)
}

func ensureDockerNetwork(name string) error {
	fmt.Printf("Ensuring Docker network %q exists...\n", name)

	// Check if the network exists
	checkCmd := exec.Command("docker", "network", "inspect", name)
	if err := checkCmd.Run(); err == nil {
		fmt.Println("Network already exists.")
		return nil
	}

	fmt.Printf("Creating Docker network %q...\n", name)
	createCmd := exec.Command("docker", "network", "create", name)
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr

	if err := createCmd.Run(); err != nil {
		return fmt.Errorf("failed to create network %q: %w", name, err)
	}

	fmt.Println("Network created successfully.")
	return nil
}
