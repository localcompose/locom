package locom

import (
	"fmt"

	"github.com/localcompose/locom/internal/hosts"
	"github.com/spf13/cobra"
)

var cmdHosts = &cobra.Command{
	Use:          "hosts",
	Short:        "Update /etc/hosts with entries from locom stage",
	Annotations: map[string]string{
		"helpdisplayorder": "40",
	},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHosts(cmd)
	},
}

func init() {
	cmdHosts.Flags().Bool("verify", false, "Check if the DNS name resolves and responds")
	rootCmd.AddCommand(cmdHosts)
}

func runHosts(cmd *cobra.Command) error {
	verify, err := cmd.Flags().GetBool("verify")
	if err != nil {
		return fmt.Errorf("failed to read verify flag: %w", err)
	}
	return hosts.Setup(verify)
}
