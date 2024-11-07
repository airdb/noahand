package oskit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// There is 2 ways to run a go program:
// 1. go run main.go
// 2. go build -o main main.go && ./main
// This function is used to check which way the program is running
// Because the diffent way will use different config path.
func IsRunGoBuild() bool {
	path := os.Args[0]

	if strings.Contains(path, "go-build") {
		log.Println("Program is running with go run")
	} else {
		log.Println("Program was built with go build")
	}

	return strings.Contains(path, "go-build")
}

func runInWindows(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return strings.TrimSpace(string(result)), err
}

func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}

func runInLinux(cmd string) (string, error) {
	log.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return strings.TrimSpace(string(result)), err
}

// 根据进程名判断进程是否运行
func CheckProRunning(serverName string) (bool, error) {
	a := `ps ux | awk '/` + serverName + `/' && !/awk/ '{print $2}'`
	pid, err := RunCommand(a)

	log.Println("xxx", err, pid)

	if err != nil {
		return false, fmt.Errorf("failed to check process running: %w", err)
	}

	return pid != "", nil
}

// 根据进程名称获取进程ID
func GetPid(serverName string) (string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return "", fmt.Errorf("failed to get process ID: %w", err)
	}

	return pid, nil
}
