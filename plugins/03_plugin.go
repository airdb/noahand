// pluginA.go
package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// Main function of the plugin.
func Plugin03() {
	log.SetFlags(log.Lshortfile)
	log.Println("Hello from Plugin 03!")

	// Create cgroup configuration
	cgroupConfig := &exec.Cmd{
		SysProcAttr: &syscall.SysProcAttr{},
	}

	cmd := exec.Command("/bin/bash", "-c", "echo $PATH && env")
	cmd.SysProcAttr = cgroupConfig.SysProcAttr
	cmd.Env = os.Environ()

	// Set standard output and error output first
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command (this will initialize cmd.Process)
	if err := cmd.Start(); err != nil {
		log.Printf("Command start failed: %v\n", err)
		return
	}

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution failed: %v\n", err)
		return
	}
}
