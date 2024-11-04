/*
Copyright © 2023 dean <dean@airdb.com>
More information, please vist https://airdb.team.
*/
package cmd

import (
	"guardhouse/coremain"

	"github.com/spf13/cobra"
)

// pluginCmd represents the plugin command.
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		coremain.RunPlugin()
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}
