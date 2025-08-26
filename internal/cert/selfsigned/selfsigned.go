package selfsigned

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Public API
//   Setup(): generates a local self-signed CA and a server cert (with SANs),
//            writes Traefik TLS config pointing to the fullchain.
//   Trust(): installs the CA into the OS trust store (curl + Chrome/Chromium on Linux,
//            System keychain on macOS, User Root on Windows). Firefox/NSS not handled yet.
//   Untrust(): removes the CA from the OS trust store.
//   Cleanup(): removes generated files (does not touch OS trust stores).

const (
	defaultCertsDir  = "./proxy/certs"
	defaultConfigDir = "./proxy/config"

	caCertName     = "selfsigned.ca.crt"
	caKeyName      = "selfsigned.ca.key"
	serverCertName = "selfsigned.server.crt"
	serverKeyName  = "selfsigned.server.key"
	fullchainName  = "selfsigned.server.fullchain.crt" // server + (optionally) intermediates/root
	traefikTLSFile = "selfsigned.yml"

	// Identity — adjust to your domain set as needed
	wildcardCN = "*.locom.self"
)

var defaultSANs = []string{
	"proxy.locom.self",
	"*.locom.self",
}

// Setup generates a CA and a server certificate (signed by that CA),
// writes PEM files with sane permissions, and creates a Traefik TLS snippet
// that references the fullchain + server key.
func Setup() error {
	if err := os.MkdirAll(defaultCertsDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(defaultConfigDir, 0o755); err != nil {
		return err
	}

	caCertPath := filepath.Join(defaultCertsDir, caCertName)
	caKeyPath := filepath.Join(defaultCertsDir, caKeyName)
	serverCertPath := filepath.Join(defaultCertsDir, serverCertName)
	serverKeyPath := filepath.Join(defaultCertsDir, serverKeyName)
	fullchainPath := filepath.Join(defaultCertsDir, fullchainName)
	traefikPath := filepath.Join(defaultConfigDir, traefikTLSFile)

	// 1) Generate CA
	caPriv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("generate CA key: %w", err)
	}
	caTpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Local Dev CA"}, CommonName: "Local Dev Root CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          mustSubjectKeyID(&caPriv.PublicKey),
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTpl, caTpl, &caPriv.PublicKey, caPriv)
	if err != nil {
		return fmt.Errorf("create CA cert: %w", err)
	}
	if err := writePEM(caCertPath, "CERTIFICATE", caDER, 0o644); err != nil {
		return err
	}
	if err := writePEM(caKeyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caPriv), 0o600); err != nil {
		return err
	}

	// 2) Generate server cert signed by CA with SANs
	srvPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate server key: %w", err)
	}
	srvTpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: wildcardCN},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     append([]string{}, defaultSANs...),
	}
	srvDER, err := x509.CreateCertificate(rand.Reader, srvTpl, caTpl, &srvPriv.PublicKey, caPriv)
	if err != nil {
		return fmt.Errorf("create server cert: %w", err)
	}
	if err := writePEM(serverCertPath, "CERTIFICATE", srvDER, 0o644); err != nil {
		return err
	}
	if err := writePEM(serverKeyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(srvPriv), 0o600); err != nil {
		return err
	}

	// 3) Fullchain (server + CA). Traefik is fine with a bundle as certFile.
	if err := concatFiles(fullchainPath, serverCertPath, caCertPath); err != nil {
		return err
	}

	// 4) Traefik dynamic TLS config snippet (paths inside the container mount)
	// Adjust mount so that host ./proxy/certs is mapped to /certs in the traefik container
	traefikYAML := fmt.Sprintf(`tls:
  certificates:
    - certFile: "/certs/%s"
      keyFile: "/certs/%s"
`, filepath.Base(fullchainPath), filepath.Base(serverKeyPath))
	if err := os.WriteFile(traefikPath, []byte(traefikYAML), 0o644); err != nil {
		return fmt.Errorf("write traefik TLS file: %w", err)
	}

	return nil
}

// Trust installs the CA into the OS trust store. Requires privileges on Linux/macOS.
func Trust() error {
	caCertPath := filepath.Join(defaultCertsDir, caCertName)
	if _, err := os.Stat(caCertPath); err != nil {
		return fmt.Errorf("CA not found: %w", err)
	}
	switch runtime.GOOS {
	case "linux":
		return linuxTrust(caCertPath)
	case "darwin":
		return darwinTrust(caCertPath)
	case "windows":
		return windowsTrust(caCertPath)
	default:
		return fmt.Errorf("trust not implemented for %s", runtime.GOOS)
	}
}

// Untrust removes the CA from the OS trust store using its fingerprint.
func Untrust() error {
	caCertPath := filepath.Join(defaultCertsDir, caCertName)
	sha, err := fileSHA1Fingerprint(caCertPath)
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "linux":
		return linuxUntrust()
	case "darwin":
		return darwinUntrust(sha)
	case "windows":
		return windowsUntrust(sha)
	default:
		return fmt.Errorf("untrust not implemented for %s", runtime.GOOS)
	}
}

// Cleanup removes generated files (does not edit trust stores)
func Cleanup() error {
	paths := []string{
		filepath.Join(defaultCertsDir, caCertName),
		filepath.Join(defaultCertsDir, caKeyName),
		filepath.Join(defaultCertsDir, serverCertName),
		filepath.Join(defaultCertsDir, serverKeyName),
		filepath.Join(defaultCertsDir, fullchainName),
		filepath.Join(defaultConfigDir, traefikTLSFile),
	}
	var errs []string
	for _, p := range paths {
		_ = os.Remove(p)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// ---- utilities ----

func writePEM(path, typ string, der []byte, mode os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: typ, Bytes: der})
}

func concatFiles(dst string, parts ...string) error {
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, p := range parts {
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if _, err := out.Write(b); err != nil {
			return err
		}
		if len(b) > 0 && b[len(b)-1] != '\n' {
			if _, err := out.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	return nil
}

func mustSubjectKeyID(pub any) []byte {
	spki, _ := x509.MarshalPKIXPublicKey(pub)
	s := sha1.Sum(spki)
	return s[:]
}

func fileSHA1Fingerprint(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return "", fmt.Errorf("no PEM in %s", path)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	fp := sha1.Sum(cert.Raw)
	return strings.ToUpper(hex.EncodeToString(fp[:])), nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
