//go:build windows

package selfsigned

import (
	"strings"
)

func trust(caCertPath string) error {
	// Add to current user Root store
	return run("certutil", "-addstore", "-user", "Root", caCertPath)
}

func untrust(sha1hex string) error {
	// Remove by SHA1 thumbprint from current user Root store
	// certutil expects hex without spaces
	return run("certutil", "-delstore", "-user", "Root", strings.ToUpper(sha1hex))
}
