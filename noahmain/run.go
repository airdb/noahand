package noahmain

import (
	"log"
	"math/rand"
	"time"

	"github.com/airdb/noah/internal/noahlib"
)

func Run() {
	for {
		// rand.Seed(900)
		maxSleepInterval := 60

		// nolint: gosec
		t := rand.Intn(maxSleepInterval)

		log.Println("sleep time seed", t)

		time.Sleep(time.Duration(t) * time.Second)

		noahlib.Heartbeat()
	}
}
