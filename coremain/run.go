package coremain

import (
	"log"
	"time"

	"guardhouse/internal/noahlib"
)

func Run() {
	// Supervisor process.
	noahlib.StartSupervisor()
	noahlib.WaitWorkerSignals()

	// Worker process.
	for {
		noahlib.DoSelfUpdate()
		log.Println("self update")
		time.Sleep(time.Minute)
	}

}
