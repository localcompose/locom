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

	return nil
}

func Trust() error {
	certPath := filepath.Join("./proxy/certs", "selfsigned.pem")
	switch runtime.GOOS {
	case "darwin":
		return darwinTrustCert(certPath)
	case "linux":
		return linuxTrustCert(certPath)
	case "windows":
		return windowsTrustCert(certPath)
	default:
		return fmt.Errorf("trust not implemented for %s", runtime.GOOS)
	}
}

func Untrust() error {
	switch runtime.GOOS {
	case "darwin":
		return darwinUntrust()
	case "linux":
		return linuxUntrust()
	case "windows":
		return windowsUntrust()
	}
	return nil
}

func Cleanup() error {
	certsDir := "./proxy/certs"
	configFile := "./proxy/config/selfsigned.yml"

	os.RemoveAll(certsDir)
	os.Remove(configFile)

	return nil
}
