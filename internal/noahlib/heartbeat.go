package noahlib

import (
	"guardhouse/internal/version"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type HostReq struct {
	IP        string `url:"ip"`
	OS        string `url:"os"`
	Hostname  string `url:"hostname"`
	Timestamp string `url:"timestamp"`
	Arch      string `url:"arch"`
	IsStart   string `url:"is_start,omitempty"`
	Username  string `url:"username"`
	Version   string `url:"version"`
}

type HostResp struct{}

// GetLocalIP returns the non loopback local IP of the host.
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func RandomHeartbeat() {
	for {
		// rand.Seed(900)
		maxSleepInterval := 60

		t := rand.Intn(maxSleepInterval)

		// t = 10
		log.Println("sleep time seed", t)

		time.Sleep(time.Duration(t) * time.Second)

		/*
			if version.GetDeployVersion() != version.GetRunningVersion() {
				SendReloadSignal()
			}
		*/

		Heartbeat()
	}
}

func Heartbeat() {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"ip":      GetLocalIP(),
			"version": version.ToString(),
		}).
		Post(GetConfigURL())
	if err != nil {
		log.Printf("Failed to send heartbeat: %v", err)

		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Heartbeat request failed with status: %d", resp.StatusCode())

		return
	}

	log.Printf("Heartbeat sent successfully")
}
