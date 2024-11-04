// pluginA.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func setupCgroup(pid int) error {
	// Define cgroup base path
	cgroupBasePath := "/sys/fs/cgroup"

	// Memory controller path
	memoryCgroupPath := filepath.Join(cgroupBasePath, "memory", "mycontainer")
	if err := os.MkdirAll(memoryCgroupPath, 0755); err != nil {
		return fmt.Errorf("failed to create memory cgroup directory: %v", err)
	}

	// Set memory limit to 100MB
	if err := os.WriteFile(filepath.Join(memoryCgroupPath, "memory.limit_in_bytes"), []byte("104857600"), 0644); err != nil {
		return fmt.Errorf("failed to set memory limit: %v", err)
	}

	// Add process to memory cgroup
	if err := os.WriteFile(filepath.Join(memoryCgroupPath, "tasks"), []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("将进程添加到内存 cgroup 失败: %v", err)
	}

	// CPU 控制器路径
	cpuCgroupPath := filepath.Join(cgroupBasePath, "cpu", "mycontainer")
	if err := os.MkdirAll(cpuCgroupPath, 0755); err != nil {
		return fmt.Errorf("创建 CPU cgroup 目录失败: %v", err)
	}

	// 设置 CPU 调度周期为 100ms
	if err := os.WriteFile(filepath.Join(cpuCgroupPath, "cpu.cfs_period_us"), []byte("100000"), 0644); err != nil {
		return fmt.Errorf("设置 CPU 调度周期失败: %v", err)
	}

	// 设置 CPU 配额为 50ms，表示 50% 的 CPU 使用率
	if err := os.WriteFile(filepath.Join(cpuCgroupPath, "cpu.cfs_quota_us"), []byte("50000"), 0644); err != nil {
		return fmt.Errorf("设置 CPU 配额失败: %v", err)
	}

	// 将进程添加到 CPU cgroup
	if err := os.WriteFile(filepath.Join(cpuCgroupPath, "tasks"), []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("将进程添加到 CPU cgroup 失败: %v", err)
	}

	return nil
}
