package oskit

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func InstallProcess() {
	tmpPath := "/sbin/noah"
	executable := "/tmp/noah_latest"

	defer os.Remove(tmpPath)
	err := exec.CommandContext(context.Background(), "/usr/bin/install", tmpPath, executable).Run()
	if err != nil {
		log.Println(err)
	}
}

func SendReloadSignal() error {
	ppid := strconv.Itoa(os.Getppid())
	err := exec.CommandContext(context.Background(), "/bin/kill", "-HUP", ppid).Run()

	return err
}
