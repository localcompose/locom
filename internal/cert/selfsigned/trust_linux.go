//go:build linux

package selfsigned

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var linuxIsOtherBrowers = false

func trust(certPath string) error {
	// 1) Install into system trust store
	dest := "/usr/local/share/ca-certificates/" + filepath.Base(certPath)
	if err := run("sudo", "cp", certPath, dest); err != nil {
		return err
	}
	if err := run("sudo", "update-ca-certificates"); err != nil {
		return err
	}

	// 2) Install into NSS DB for Chrome
	if err := linuxChromeAddToNSSDB(certPath); err != nil {
		fmt.Println("⚠ CHrome NSS DB update skipped:", err)
		fmt.Println("Install libnss3-tools and rerun Trust() to fix Chrome/Chromium trust.")
	}

	return nil
}

func untrust(_ string) error {
	// 1) Remove from system CA store
	cand := "/usr/local/share/ca-certificates/" + caCertName
	_ = run("sudo", "rm", "-f", cand)
	_ = run("sudo", "update-ca-certificates")

	// 2) Remove from NSS DB
	if err := linuxChromeRemoveFromNSSDB(); err != nil {
		fmt.Println("⚠ Could not remove from NSS DB:", err)
	}

	return nil
}

func linuxChromeAddToNSSDB(certPath string) error {
	// Check if certutil exists
	if _, err := exec.LookPath("certutil"); err != nil {
		return fmt.Errorf("certutil not found")
	}

	nssdb := filepath.Join(os.Getenv("HOME"), ".pki", "nssdb")
	if err := os.MkdirAll(nssdb, 0o755); err != nil {
		return fmt.Errorf("failed to create NSS DB folder: %w", err)
	}

	cmd := exec.Command("certutil", "-d", "sql:"+nssdb, "-A", "-t", "C,,",
		"-n", "locom-selfsigned", "-i", certPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add cert to NSS DB: %v, output: %s", err, string(output))
	}

	return nil
}

func linuxChromeRemoveFromNSSDB() error {
	if _, err := exec.LookPath("certutil"); err != nil {
		return fmt.Errorf("certutil not found")
	}

	nssdb := filepath.Join(os.Getenv("HOME"), ".pki", "nssdb")
	cmd := exec.Command("certutil", "-d", "sql:"+nssdb, "-D", "-n", "locom-selfsigned")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove cert from NSS DB: %v, output: %s", err, string(output))
	}

	return nil
}
