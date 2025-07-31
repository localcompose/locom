package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"net"
	"time"
	"syscall"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var cmdHosts = &cobra.Command{
	Use:   "hosts",
	Short: "Update /etc/hosts with entries from locom stage",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHosts(cmd)
	},
}

func init() {
	cmdHosts.Flags().Bool("verify", false, "Check if the DNS name resolves and responds")
	rootCmd.AddCommand(cmdHosts)
}

func runHosts(cmd *cobra.Command) error {
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

	cpCmd := exec.Command("sudo", "cp", tmpHosts.Name(), hostsPath)
	cpCmd.Stdin = os.Stdin
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("updating /etc/hosts with sudo: %w", err)
	}

	statePath := ".locom/hosts"
	if err := ioutil.WriteFile(statePath, []byte(fmt.Sprintf("%s\n%s\n%s\n", beginMarker, entry, endMarker)), 0644); err != nil {
		return fmt.Errorf("writing state to .locom/hosts: %w", err)
	}

	fmt.Println("✅ Hosts file updated with locom stage entries.")

	verify, _ := cmd.Flags().GetBool("verify")
	if verify {
		fqdn := "proxy" + suffix
		if err := verifyHost(address, fqdn); err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}
	}

	return nil
}

func verifyHost(expectedAddr, fqdn string) error {
	fmt.Printf("🔍 Verifying DNS resolution for %s...\n", fqdn)

	ips, err := net.LookupHost(fqdn)
	if err != nil {
		return fmt.Errorf("DNS resolution failed for %s: %w", fqdn, err)
	}

	match := false
	for _, ip := range ips {
		if ip == expectedAddr {
			match = true
			break
		}
	}
	if !match {
		return fmt.Errorf("DNS %s resolved to %v, expected %s", fqdn, ips, expectedAddr)
	}

	fmt.Printf("✅ DNS resolution successful: %s → %s\n", fqdn, expectedAddr)

	// Optional: attempt TCP connection to verify routing
	addr := net.JoinHostPort(fqdn, "80")
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok {
			if errors.Is(opErr.Err, syscall.ECONNREFUSED) {
				fmt.Printf("⚠️ TCP Connection to %s at port 80 refused (no service), but DNS resolution succeeded.\n", addr)
				return nil
			}
		}
		fmt.Printf("⚠️ TCP connection failed to %s: %v", addr, err)
		return nil
	}
	fmt.Printf("✅ TCP connection successful to %s\n", addr)
	return nil
}
