package locom

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/template"

	"github.com/spf13/cobra"
)

var main struct {
	Name    string
	Version string
}

var rootCmd = &cobra.Command{
	Use:   "locom",
	Short: "locom manages a local stage of Docker Compose stacks",
	Long:  `locom is a CLI tool for managing local Docker Compose stacks in a minimal, offline-friendly way.`,
}

func NewRootCmd() *cobra.Command {
	return rootCmd
}

// Execute runs the root command
func Execute(name, version string) {
	main.Name = name
	main.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register custom template function globally
	cobra.AddTemplateFuncs(template.FuncMap{
		"orderedCommands": orderedCommands,
	})

	rootCmd.SetUsageTemplate(`{{if .HasExample}}{{.Example}}

{{end}}{{if .Short}}{{.Short}}{{end}}

Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}} [command]{{end}}

Available Commands:{{range orderedCommands .}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}

{{if .HasAvailableLocalFlags}}Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableSubCommands}}Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

// Sort commands by annotation
func orderedCommands(cmd *cobra.Command) []*cobra.Command {
	orderhelpcompletion()
	cmds := append([]*cobra.Command(nil), cmd.Commands()...)
	sort.SliceStable(cmds, func(i, j int) bool {
		oi := orderValue(cmds[i])
		oj := orderValue(cmds[j])
		if oi == oj {
			return cmds[i].Name() < cmds[j].Name()
		}
		return oi < oj
	})
	return cmds
}

func orderValue(c *cobra.Command) int {
	if v, ok := c.Annotations["helpdisplayorder"]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 9999
}

func orderhelpcompletion() {
	if c := rootCmd.Commands(); len(c) > 0 {
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "completion" {
				if cmd.Annotations == nil {
					cmd.Annotations = map[string]string{}
				}
				cmd.Annotations["helpdisplayorder"] = "2"
			}
			if cmd.Name() == "help" {
				if cmd.Annotations == nil {
					cmd.Annotations = map[string]string{}
				}
				cmd.Annotations["helpdisplayorder"] = "1"
			}
		}
	}
}
