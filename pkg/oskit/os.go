package oskit

import (
	"fmt"
	"guardhouse/pkg/configkit"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	InstallTypeForce     = "force"
	InstallTypeOverwrite = "overwrite"
	InstallTypeSkip      = "skip"
	InstallTypeUpdate    = "update"
)

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

func InstallTarball() {
	// 1. Download the tarball
	// 2. Calculate the MD5 checksum of the tarball
	// 3. Extract the tarball
	// 4. Install the modules by directly copying the files

	// Download a file
	surl := "http://oracle.airdb.host:8000/noah_latest.tgz"
	dpath := "/tmp/noah_latest.tgz"

	err := ResumableDownload(surl, dpath)
	if err != nil {
		log.Println("Download failed")
		return
	}

	log.Println("Download successfully")

	// get file md5
	md5, err := GetFileMD5(dpath)
	if err != nil {
		log.Println("Failed to get file md5")
		return
	}
	unzipPath := path.Join("/tmp/noah", md5)

	err = os.MkdirAll(unzipPath, 0o755)
	if err != nil {
		log.Println("Failed to create directory", unzipPath)
		return
	}

	defer os.Remove(dpath)
	ExtractTarball(dpath, unzipPath)
	log.Println("Unzipped to", unzipPath)

	// The runtime process cannot install directly.
	modulesPath := path.Join(configkit.ModulesDir, md5)
	InstallDirectory(unzipPath, modulesPath)
}
