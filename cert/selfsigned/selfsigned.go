package selfsigned

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Setup() error {
	certsDir := "./proxy/certs"
	configDir := "./proxy/config"
	os.MkdirAll(certsDir, 0755)
	os.Chmod(certsDir, 0755)
	os.MkdirAll(configDir, 0755)
	os.Chmod(configDir, 0755)

	crtFile := filepath.Join(certsDir, "selfsigned.pem")
	keyFile := filepath.Join(certsDir, "selfsigned-key.pem")
	cn := "*.locom.self"

	cmd := exec.Command("openssl", "req", "-x509", "-nodes", "-days", "825",
		"-newkey", "rsa:2048",
		"-subj", "/CN="+cn,
		"-keyout", keyFile,
		"-out", crtFile,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("openssl failed: %w", err)
	}

	tlsConf := filepath.Join(configDir, "selfsigned.yml")
	yml := `tls:
  certificates:
    - certFile: "/certs/selfsigned.pem"
      keyFile: "/certs/selfsigned-key.pem"
`
	if err := os.WriteFile(tlsConf, []byte(yml), 0644); err != nil {
		return err
	}

	// return trustCert(crtFile)
	return nil
}

func Trust() error {
	certsDir := "./proxy/certs"

	crtFile := filepath.Join(certsDir, "selfsigned.pem")
	return trustCert(crtFile)
}

func trustCert(certPath string) error {
	switch runtime.GOOS {
	case "darwin":
		// macOS system trust store
		return exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot",
			"-k", "/Library/Keychains/System.keychain", certPath).Run()

	case "linux":
		// 1. Install into system trust store (OpenSSL/GnuTLS consumers like curl, git)
		if err := exec.Command("sudo", "cp", certPath, "/usr/local/share/ca-certificates/locom-selfsigned.crt").Run(); err != nil {
			return err
		}
		if err := exec.Command("sudo", "update-ca-certificates").Run(); err != nil {
			return err
		}

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

	case "windows":
		// Windows user root store
		return exec.Command("certutil", "-addstore", "-user", "Root", certPath).Run()

	default:
		return fmt.Errorf("trust not implemented for %s", runtime.GOOS)
	}
}

func Untrust() error {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("sudo", "security", "delete-certificate", "-c", "*.locom.self").Run()
	case "linux":
		exec.Command("sudo", "rm", "/usr/local/share/ca-certificates/locom-selfsigned.crt").Run()
		exec.Command("sudo", "update-ca-certificates").Run()
	case "windows":
		exec.Command("certutil", "-delstore", "-user", "Root", "*.locom.self").Run()
	}

	return nil
}

func Cleanup() error {
	certsDir := "./proxy/certs"
	configFile := "./proxy/config/selfsigned.yml"
	// crtFile := filepath.Join(certsDir, "selfsigned.pem")

	// switch runtime.GOOS {
	// case "darwin":
	// 	exec.Command("sudo", "security", "delete-certificate", "-c", "*.locom.self").Run()
	// case "linux":
	// 	exec.Command("sudo", "rm", "/usr/local/share/ca-certificates/locom-selfsigned.crt").Run()
	// 	exec.Command("sudo", "update-ca-certificates").Run()
	// case "windows":
	// 	exec.Command("certutil", "-delstore", "-user", "Root", "*.locom.self").Run()
	// }

	os.RemoveAll(certsDir)
	os.Remove(configFile)
	return nil
}
