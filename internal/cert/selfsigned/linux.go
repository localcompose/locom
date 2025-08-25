package selfsigned

import (
	"os"
	"os/exec"
	"path/filepath"
)

func linuxTrustCert(certPath string) error {
	// 1. Install into system trust store (OpenSSL/GnuTLS consumers like curl, git)
	if err := exec.Command("sudo", "cp", certPath, "/usr/local/share/ca-certificates/locom-selfsigned.crt").Run(); err != nil {
		return err
	}
	if err := exec.Command("sudo", "update-ca-certificates").Run(); err != nil {
		return err
	}

	if err := linuxNss(certPath); err != nil {
		return err
	}

	return nil
}

func linuxNss(certPath string) error {
	isNss := false
	if isNss {
		// 2. Also add to user NSS DB (for Chromium/Firefox/Puppeteer)
		// Ensure nss-tools installed: sudo apt install libnss3-tools
		// Create DB if missing
		if err := exec.Command("mkdir", "-p", filepath.Join(os.Getenv("HOME"), ".pki", "nssdb")).Run(); err != nil {
			return err
		}
		nssdb := "sql:" + filepath.Join(os.Getenv("HOME"), ".pki", "nssdb")
		if err := exec.Command("certutil", "-d", nssdb, "-A", "-t", "C,,", "-n",
			"locom-selfsigned", "-i", certPath).Run(); err != nil {
			return err
		}
	}
	return nil
}

func linuxUntrust() error {
	exec.Command("sudo", "rm", "/usr/local/share/ca-certificates/locom-selfsigned.crt").Run()
	exec.Command("sudo", "update-ca-certificates").Run()

	return nil
}
