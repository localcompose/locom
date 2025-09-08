package locom

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run:   runVersion,
	Annotations: map[string]string{
		"helpdisplayorder": "10",
	},
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s %s %s\n", main.Name, cmd.Use, main.Version)
}

func init() {
	rootCmd.AddCommand(cmdVersion)
}
