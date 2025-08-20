package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var cmdHosts = &cobra.Command{
	Use:          "hosts",
	Short:        "Update /etc/hosts with entries from locom stage",
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

	// hostsPath := `C:\Windows\System32\drivers\etc\hosts` // "C:\\Program Files\\Git\\etc\\hosts" // "/etc/hosts"
	hostsPath := `C:\Windows\System32\drivers\etc\hosts`
	if _, err := os.Stat(hostsPath); err != nil {
		// if System32 is redirected, fall back to Sysnative
		hostsPath = `C:\Windows\Sysnative\drivers\etc\hosts`
	}

	tmpHosts, err := os.CreateTemp("", "hosts.*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpHosts.Name())

	hostsContent, err := ioutil.ReadFile(hostsPath)
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

	f, err := os.OpenFile(hostsPath, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
			return rerunAsAdmin(err, tryRunAsAdmin)
		}
		return fmt.Errorf("opening hosts file for write: %w", err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(updated)); err != nil {
		if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
			return rerunAsAdmin(err, tryRunAsAdmin)
		}
		return fmt.Errorf("writing hosts file: %w", err)
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

var tryRunAsAdmin = true

func rerunAsAdmin(origErr error, tryRunAsAdmin bool) error {
	exe, _ := os.Executable()
	args := strings.Join(os.Args[1:], " ")

	if tryRunAsAdmin {
		// silently re-run as admin
		if err := runAsAdmin(exe, args); err != nil {
			return fmt.Errorf("access denied, tried to elevate: %w", err)
		}
		os.Exit(0) // stop current process, elevated one continues
	}

	// suggest the command instead
	fmt.Printf("\nAccess denied while writing hosts file.\n")
	fmt.Printf("You can retry with Administrator rights:\n\n")
	fmt.Printf("  runas /user:Administrator \"%s %s\"\n\n", exe, args)
	return origErr
}

func runAsAdmin(exePath, args string) error {
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	argPtr, _ := syscall.UTF16PtrFromString(args)
	cwdPtr, _ := syscall.UTF16PtrFromString("")

	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShellExecute := shell32.NewProc("ShellExecuteW")

	r, _, _ := procShellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argPtr)),
		uintptr(unsafe.Pointer(cwdPtr)),
		1, // SW_NORMAL
	)
	if r <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", r)
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
