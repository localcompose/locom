package cmd

import (
	"github.com/localcompose/locom/cert/selfsigned"
	"github.com/spf13/cobra"
)

func init() {
	cmdCert.AddCommand(cmdSelfSigned)
	cmdSelfSigned.AddCommand(cmdSelfSignedUninstall)

	rootCmd.AddCommand(cmdCert)
}

var cmdCert = &cobra.Command{
	Use:   "cert",
	Short: "Manage certificates for locom",
}

var cmdSelfSigned = &cobra.Command{
	Use:   "selfsigned",
	Short: "Generate and trust a self-signed certificate for .locom.self",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfsigned.Setup()
	},
}

var cmdSelfSignedUninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove self-signed cert and its trust config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfsigned.Cleanup()
	},
}
