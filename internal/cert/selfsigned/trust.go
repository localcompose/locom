package selfsigned

import (
	"fmt"
	"os"
	"path/filepath"
)

// TrustSetup installs the CA into the OS trust store. Requires privileges on Linux/macOS.
func TrustSetup() error {
	caCertPath := filepath.Join(defaultCertsDir, caCertName)
	if _, err := os.Stat(caCertPath); err != nil {
		return fmt.Errorf("CA not found: %w", err)
	}

	return trust(caCertPath)
}

// TrustCleanup removes the CA from the OS trust store using its fingerprint.
func TrustCleanup() error {
	caCertPath := filepath.Join(defaultCertsDir, caCertName)
	sha, err := fileSHA1Fingerprint(caCertPath)
	if err != nil {
		return err
	}

	return untrust(sha)
}
