//go:build darwin || linux

package hosts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func getHostsPath() string {
	return "/etc/hosts"
}

func updateHosts(updatedContent, hostsPath string) error {
	tryElevated := true

	tmpHosts, err := os.CreateTemp(filepath.Dir(hostsPath), "hosts.*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpHostsPath := tmpHosts.Name()
	defer os.Remove(tmpHostsPath)

	if err := os.WriteFile(tmpHostsPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("writing temp hosts file: %w", err)
	}

	// Try direct copy first
	err = copyFile(tmpHostsPath, hostsPath)
	if !os.IsPermission(err) {
		return err
	}

	if tryElevated {
		// Permission denied → retry with sudo tee
		return unixCopyWithInteractiveElevation(tmpHostsPath, hostsPath)
	}

	return err
}

func unixCopyWithInteractiveElevation(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open temp hosts: %w", err)
	}
	defer src.Close()

	cmd := exec.Command("sudo", "tee", dstPath)
	cmd.Stdin = src
	cmd.Stdout = os.Stdout // so user sees “Password:” prompt
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sudo tee %s failed: %w", dstPath, err)
	}
	return nil
}
