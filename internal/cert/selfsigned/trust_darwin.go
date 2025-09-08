//go:build darwin

package selfsigned

import (
	"strings"
)

func trust(caCertPath string) error {
	// Add to System keychain as a trusted root (requires sudo)
	return run("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain", caCertPath)
}

func untrust(sha1hex string) error {
	// Remove by fingerprint from System keychain
	return run("sudo", "security", "delete-certificate", "-Z", strings.ToUpper(sha1hex), "/Library/Keychains/System.keychain")
}
