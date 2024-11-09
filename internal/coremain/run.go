package coremain

import (
	"fmt"
	"guardhouse/internal/admin"
	"guardhouse/internal/noahlib"
	"guardhouse/pkg/configkit"
	"log"
	"os"
	"os/user"
	"time"
)

func Run() {
	if !IsRuuserRoot() {
		return
	}

	configkit.InitConfig()

	// Supervisor process.
	noahlib.StartSupervisor()

	go admin.RunServer()
	time.Sleep(time.Second * 3)

	// Worker process.
	go func() {
		log.SetPrefix(fmt.Sprintf("PID-%d ", os.Getpid()))

		executable, _ := os.Executable()
		log.Println("noah worker, executable, ", executable)
		// noahlib.MonitorMasterProcess(executable)

		for {
			log.Println("plugin management")
			RunPlugin()
			// noahlib.DoSelfUpdate()
			// log.Println("self update")
			time.Sleep(time.Minute)
		}
	}()

	noahlib.WaitWorkerSignals()
}

func IsRuuserRoot() bool {
	// Check run user
	runUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
		return false
	}

	if runUser.Uid == "0" {
		return true
	}

	log.Fatalf("current user: %s, must run as root!\n", runUser.Username)

	return false
}
