/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"guardhouse/internal/noahlib"
	"guardhouse/pkg/version"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "noah",
	Short: "Noah Server Management System",
	Long: `Noah is a server management system that helps you manage your servers efficiently.

Get Started:
- Step 1: run "noah init" to generate a systemd service file.
- Step 2: run "sudo systemctl daemon-reload" to reload the systemd service.
- Step 3: run "udo systemctl start noah" to start the noah service.
- Step 4: run "sudo systemctl enable noah" to enable the noah service on boot.
- Step 5: run "sudo journalctl -u noah -f" to view the noah service logs.

For more information, please visit https://aid.run`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// updateCmd represents the update command.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		noahlib.DoSelfUpdate()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = version.ToString()
	rootCmd.AddCommand(updateCmd)
	initRun()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
