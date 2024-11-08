/*
Copyright © 2023 dean <dean@airdb.com>
More information, please vist https://airdb.team.
*/
package cmd

import (
	"fmt"
	"guardhouse/pkg/configkit"
	"guardhouse/pkg/service"
	"log"
	"runtime"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitSystemdService()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func InitSystemdService() {
	serviceFilePath := configkit.SystemdFilepath

	if runtime.GOOS != "linux" {
		log.Println("init systemd service is only supported on Linux.")
		serviceFilePath = "/tmp/noah.service"
	}

	config := service.ServiceConfig{
		Description:      "Noah Server Management System",
		Documentation:    "https://aid.run",
		ExecStart:        "/opt/noah/noah run",
		User:             "root",
		Group:            "root",
		WorkingDirectory: "/opt/noah",
		Environment:      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin",
	}

	if err := service.GenerateServiceFile(config, serviceFilePath); err != nil {
		log.Fatalf("Failed to generate service file: %v", err)
	}

	fmt.Printf("Successfully generated service file at: %s\n", serviceFilePath)
}
