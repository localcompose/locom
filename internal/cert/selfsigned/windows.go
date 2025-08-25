package selfsigned

import (
	"os/exec"
)

func windowsTrustCert(certPath string) error {
	// Windows user root stor
	return exec.Command("certutil", "-addstore", "-user", "Root", certPath).Run()
}

func windowsUntrust() error {
	exec.Command("certutil", "-delstore", "-user", "Root", "*.locom.self").Run()

	return nil
}
