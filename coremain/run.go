package coremain

import (
	"log"
	"time"

	"guardhouse/coremain/admin"
	"guardhouse/internal/noahlib"
)

func Run() {
	// Supervisor process.
	noahlib.StartSupervisor()
	go admin.RunServer()
	noahlib.WaitWorkerSignals()

	// Worker process.
	for {
		noahlib.DoSelfUpdate()
		log.Println("self update")
		time.Sleep(time.Minute)
	}

}
