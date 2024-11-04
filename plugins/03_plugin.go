// pluginA.go
package main

import (
	"log"
	"os"
	"os/exec"
)

// Main function of the plugin.
func Plugin03() {
	log.SetFlags(log.Lshortfile)
	log.Println("Hello from Plugin 03!")

	cmd := exec.Command("/bin/bash", "-c", "echo $PATH && ps && env")
	cmd.Env = os.Environ() // Inherit parent process environment variables

	// Set standard output and error output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute command
	err := cmd.Run()
	if err != nil {
		log.Printf("Command execution failed: %v\n", err)
		return
	}
}
