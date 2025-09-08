package hosts

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

func Setup(verify bool) error {
	configPath := ".locom/locom.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New("this folder does not contain locom stage configuration")
	}

	// Read YAML
	content, err := os.ReadFile(configPath)
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
		return errors.New("missing required fields in locom.yml (stage.network.bind.address or stage.network.dns.suffix)")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current dir: %w", err)
	}
	stageName := filepath.Base(cwd)

	beginMarker := fmt.Sprintf("# >>> locom %s loopback apps >>>", stageName)
	endMarker := fmt.Sprintf("# <<< locom %s loopback apps <<<", stageName)
	entry := fmt.Sprintf("%s proxy%s", address, suffix)

	hostsPath := getHostsPath()
	hostsContent, err := os.ReadFile(hostsPath)
	if err != nil {
		return fmt.Errorf("reading /etc/hosts: %w", err)
	}

	sep := "\n"
	if strings.Contains(string(hostsContent), "\r\n") {
		sep = "\r\n"
	}

	lines := strings.Split(string(hostsContent), sep)
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

	updated := strings.Join(newLines, sep) + sep

	if err := updateHosts(updated, hostsPath); err != nil {
		return err
	}

	statePath := ".locom/hosts"
	if err := os.WriteFile(statePath, []byte(fmt.Sprintf("%s\n%s\n%s\n", beginMarker, entry, endMarker)), 0644); err != nil {
		return fmt.Errorf("writing state to .locom/hosts: %w", err)
	}

	fmt.Println("✅ Hosts file updated with locom stage entries.")

	if verify {
		fqdn := "proxy" + suffix
		if err := verifyHost(address, fqdn); err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open temp hosts: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
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
