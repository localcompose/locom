package selfsigned

import (
	"strings"
)

func darwinTrust(caCertPath string) error {
	// Add to System keychain as a trusted root (requires sudo)
	return run("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain", caCertPath)
}

func darwinUntrust(sha1hex string) error {
	// Remove by fingerprint from System keychain
	return run("sudo", "security", "delete-certificate", "-Z", strings.ToUpper(sha1hex), "/Library/Keychains/System.keychain")
}
