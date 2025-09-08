package locom

import (
	"github.com/localcompose/locom/internal/stage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdProxy)
}

var cmdProxy = &cobra.Command{
	Use:   "proxy",
	Short: "Create a default docker-compose configuration with Traefik proxy",
	Annotations: map[string]string{
		"helpdisplayorder": "50",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "proxy"
		return stage.GenerateProxyComposeFiles(".locom/locom.yml", target)
	},
}
