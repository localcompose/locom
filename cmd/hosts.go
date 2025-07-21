package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdHosts)
}

var cmdHosts = &cobra.Command{
	Use:   "hosts",
	Short: "Update /etc/hosts with entries from locom stage",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHosts()
	},
}

func runHosts() error {
	configPath := ".locom/locom.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New("This folder does not contain locom stage configuration.")
	}

	// Read YAML
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading locom.yml: %w", err)
	}

	var parsed struct {
		Stage struct {
			Network struct {
				Bind struct {
					Address string `yaml:"address"`
				} `yaml:"bind"`
				DNS struct {
					Suffix string `yaml:"suffix"`
				} `yaml:"dns"`
			} `yaml:"network"`
		} `yaml:"stage"`
	}
	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	address := parsed.Stage.Network.Bind.Address
	suffix := parsed.Stage.Network.DNS.Suffix
	if address == "" || suffix == "" {
		return errors.New("Missing required fields in locom.yml (stage.network.bind.address or stage.network.dns.suffix)")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current dir: %w", err)
	}
	project := filepath.Base(cwd)

	beginMarker := fmt.Sprintf("# >>> locom %s loopback apps >>>", project)
	endMarker := fmt.Sprintf("# <<< locom %s loopback apps <<<", project)
	entry := fmt.Sprintf("%s proxy%s", address, suffix)

	hostsPath := "/etc/hosts"
	tmpHosts, err := os.CreateTemp("", "hosts.*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpHosts.Name())

	hostsContent, err := ioutil.ReadFile(hostsPath)
	if err != nil {
		return fmt.Errorf("reading /etc/hosts: %w", err)
	}

	lines := strings.Split(string(hostsContent), "\n")
	inBlock := false
	var newLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) == beginMarker {
			inBlock = true
			continue
		}
		if strings.TrimSpace(line) == endMarker {
			inBlock = false
			continue
		}
		if !inBlock {
			newLines = append(newLines, line)
		}
	}

	newLines = append(newLines,
		beginMarker,
		entry,
		endMarker,
	)

	updated := strings.Join(newLines, "\n") + "\n"
	if err := ioutil.WriteFile(tmpHosts.Name(), []byte(updated), 0644); err != nil {
		return fmt.Errorf("writing temp hosts file: %w", err)
	}

	cmd := exec.Command("sudo", "cp", tmpHosts.Name(), hostsPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("updating /etc/hosts with sudo: %w", err)
	}

	statePath := ".locom/hosts"
	if err := ioutil.WriteFile(statePath, []byte(fmt.Sprintf("%s\n%s\n%s\n", beginMarker, entry, endMarker)), 0644); err != nil {
		return fmt.Errorf("writing state to .locom/hosts: %w", err)
	}

	fmt.Println("✅ Hosts file updated with locom stage entries.")
	return nil
}
