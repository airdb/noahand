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

func DownloadAndInstallTarball() error {
	// 1. Download the tarball
	// 2. Calculate the MD5 checksum of the tarball
	// 3. Extract the tarball
	// 4. Install the modules by directly copying the files

	// Download a file
	surl := "http://oracle.airdb.host:8000/noah_latest.tgz"
	dpath := "/tmp/noah_latest.tgz"

	// 1. Get the md5 value from the md5 file.
	md5, err := DownloadMd5File(surl + ".md5")
	if err != nil {
		log.Println("Failed to download md5 file", err)

		return err
	}

	// 2. check if the `modules/{md5}` directory exists.
	if IsDirNotEmpty(configkit.ModulesDir + "/" + md5) {
		log.Println("modules already installed, skip installation, module md5:", md5)
		return nil
	}

	// 3. Ensure the download file is not empty.

	filesize, err := GetDownloadFileSize(surl)
	if err != nil {
		log.Println("Failed to get download file size", err)
		return err
	}

	if filesize <= 0 {
		log.Println("Download file size is 0")
		return fmt.Errorf("download file size is 0")
	}

	// 4. Download the tarball.
	err = ResumableDownload(surl, dpath)
	if err != nil {
		log.Println("Download failed", err)

		return err
	}

	// 5. Calculate the MD5 checksum of the tarball.
	md5File, err := GetFileMD5(dpath)
	if err != nil {
		log.Println("Failed to get local file md5", err)
		return err
	}

	// 6. Compare the md5 value from the md5 file and the local file.
	if md5File != md5 {
		log.Println("MD5 value does not match")
		return fmt.Errorf("md5 value does not match")
	}
	log.Println("Download successfully")

	// 7. Extract the tarball.
	unzipPath := path.Join("/tmp/noah", md5)

	ExtractTarball(dpath, unzipPath)
	log.Println("Unzipped to", configkit.ModulesDir+"/"+md5)

	defer os.Remove(dpath)
	ExtractTarball(dpath, unzipPath)
	log.Println("Unzipped to", unzipPath)

	// 8. Install the modules by directly copying the files.
	modulesPath := path.Join(configkit.ModulesDir, md5)

	if IsDirNotEmpty(modulesPath) {
		log.Println("Module already exists, skip installation, module md5:", md5)
		// InstallTypeSkip

		return nil
	}

	InstallDirectory(unzipPath, modulesPath)
	return nil
}

// check directory exist and not empty.
func IsDirNotEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == nil
}
