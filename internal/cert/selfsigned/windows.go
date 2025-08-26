package selfsigned

import (
	"strings"
)

func windowsTrust(caCertPath string) error {
	// Add to current user Root store
	return run("certutil", "-addstore", "-user", "Root", caCertPath)
}

func windowsUntrust(sha1hex string) error {
	// Remove by SHA1 thumbprint from current user Root store
	// certutil expects hex without spaces
	return run("certutil", "-delstore", "-user", "Root", strings.ToUpper(sha1hex))
}
