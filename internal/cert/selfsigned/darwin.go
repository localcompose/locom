package selfsigned

import (
	"os/exec"
)

func darwinTrustCert(certPath string) error {
	// macOS system trust store
	return exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain", certPath).Run()
}

func darwinUntrust() error {
	exec.Command("sudo", "security", "delete-certificate", "-c", "*.locom.self").Run()

	return nil
}
