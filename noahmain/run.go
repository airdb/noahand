package noahmain

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/airdb/noah/internal/noahlib"
)

func Run() {
	go Heartbeat()

	for {
		fmt.Println("run plugin...")
		time.Sleep(time.Minute)
	}
}

func Heartbeat() {
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
