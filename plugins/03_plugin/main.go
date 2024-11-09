// pluginA.go
package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Plugin03 struct{}

var Plugin = Plugin03{}

func (p Plugin03) GetInfo() {
	log.Println("Plugin03 version: v0.1.0")
}

func (p Plugin03) Init() {
	log.Println("Init from Plugin 03!")
}

func (p Plugin03) Start() {
	log.Println("Hello from Plugin 03!")
	p.Execute()
}

func (p Plugin03) Stop() {
	log.Println("Stop from Plugin 03!")
}

// Main function of the plugin.
func (p Plugin03) Execute() {
	// Create cgroup configuration
	cgroupConfig := &exec.Cmd{
		SysProcAttr: &syscall.SysProcAttr{},
	}

	cmd := exec.Command("/bin/bash", "-c", "echo plugin03 && echo $HOME")
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
