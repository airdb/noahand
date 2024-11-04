// pluginA.go
package main

import (
	"fmt"
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

	// Now we can safely access cmd.Process.Pid
	if err := setupCgroup(cmd.Process.Pid); err != nil {
		log.Printf("Failed to setup cgroup: %v\n", err)
		return
	}

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution failed: %v\n", err)
		return
	}
}

// Setup cgroup limits
func setupCgroup(pid int) error {
	// Create cgroup
	cgroupPath := "/sys/fs/cgroup/memory/mycontainer"
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return err
	}

	// Set memory limit (e.g., 100MB)
	if err := os.WriteFile(cgroupPath+"/memory.limit_in_bytes", []byte("104857600"), 0644); err != nil {
		return err
	}

	// Set CPU limit (e.g., 50% CPU)
	if err := os.WriteFile(cgroupPath+"/cpu.cfs_quota_us", []byte("1000"), 0644); err != nil {
		return err
	}

	// Add process to cgroup
	return os.WriteFile(cgroupPath+"/tasks", []byte(fmt.Sprintf("%d", pid)), 0644)
}
