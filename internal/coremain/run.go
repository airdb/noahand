package coremain

import (
	"fmt"
	"noahand/internal/admin"
	"noahand/internal/noahlib"
	"noahand/pkg/configkit"
	"log"
	"os"
	"os/user"
	"time"
)

func Run() {
	MustRunAsRoot()

	configkit.InitConfig()

	// Supervisor (master) blocks in here forever; the worker returns and
	// continues below.
	noahlib.StartSupervisor()

	go admin.RunServer()

	// Worker process.
	go func() {
		log.SetPrefix(fmt.Sprintf("PID-%d ", os.Getpid()))

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

func MustRunAsRoot() {
	runUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	if runUser.Uid != "0" {
		log.Fatalf("current user: %s, must run as root!", runUser.Username)
	}
}
