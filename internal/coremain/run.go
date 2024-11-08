package coremain

import (
	"fmt"
	"guardhouse/internal/admin"
	"guardhouse/internal/noahlib"
	"guardhouse/pkg/configkit"
	"log"
	"os"
	"time"
)

func Run() {
	log.SetFlags(log.Lshortfile)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.SetPrefix(fmt.Sprintf("%s PID-%d ", timestamp, os.Getpid()))

	configkit.InitConfig()

	// Supervisor process.
	noahlib.StartSupervisor()

	go admin.RunServer()

	// Worker process.
	go func() {
		log.SetFlags(log.Lshortfile)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		log.SetPrefix(fmt.Sprintf("%s PID-%d ", timestamp, os.Getpid()))

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
