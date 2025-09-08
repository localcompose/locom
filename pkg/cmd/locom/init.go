package locom

import (
	"github.com/localcompose/locom/internal/stage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdInit)
}

var cmdInit = &cobra.Command{
	Use:   "init [folder]",
	Short: "Initialize a new locom stage in the specified folder",
	Long: `Creates a .locom directory and a default locom.yml config file 
inside the given folder. Fails if the folder is not empty.`,
	Annotations: map[string]string{
		"helpdisplayorder": "20",
	},
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) == 1 {
			target = args[0]
		}
		return stage.Init(target)
	},
}
