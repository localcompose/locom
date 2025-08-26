package locom

import (
	"github.com/localcompose/locom/internal/cert/selfsigned"
	"github.com/spf13/cobra"
)

func init() {
	cmdCert.AddCommand(cmdSelfSigned)
	cmdSelfSigned.AddCommand(cmdSelfSignedSetup)
	cmdSelfSigned.AddCommand(cmdSelfSignedTrust)
	cmdSelfSigned.AddCommand(cmdSelfSignedUntrust)
	cmdSelfSigned.AddCommand(cmdSelfSignedCleanup)

	rootCmd.AddCommand(cmdCert)
}

var cmdCert = &cobra.Command{
	Use:   "cert",
	Short: "Manage certificates for locom",
}

var cmdSelfSigned = &cobra.Command{
	Use:   "selfsigned",
	Short: "Generate a self-signed certificate for .locom.self",
}

var cmdSelfSignedSetup = &cobra.Command{
	Use:   "setup",
	Short: "Generate a self-signed certificate for .locom.self",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfsigned.Setup()
	},
}

var cmdSelfSignedTrust = &cobra.Command{
	Use:   "trust",
	Short: "trust the self-signed certificate for .locom.self",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfsigned.Trust()
	},
}

// Untrust - remove the certificate from all trust stores
var cmdSelfSignedUntrust = &cobra.Command{
	Use:   "untrust",
	Short: "Remove/unregister the self-signed certificate from all trust stores",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfsigned.Untrust()
	},
}

var cmdSelfSignedCleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove self-signed cert and its trust config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfsigned.Cleanup()
	},
}
