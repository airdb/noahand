package coremain

import (
	"guardhouse/internal/admin"
	"guardhouse/internal/noahlib"
	"guardhouse/pkg/configkit"
	"log"
	"time"
)

func Run() {
	configkit.InitConfig()

	// Supervisor process.
	noahlib.StartSupervisor()

	go admin.RunServer()

	// Worker process.
	go func() {
		for {
			RunPlugin()
			log.Println("plugin run")
			// noahlib.DoSelfUpdate()
			// log.Println("self update")
			time.Sleep(time.Minute)
		}
	}()

	noahlib.WaitWorkerSignals()
}
