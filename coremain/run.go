package coremain

import (
	"guardhouse/internal/admin"
	"guardhouse/internal/noahlib"
	"log"
	"time"
)

func Run() {
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
