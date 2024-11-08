package noahlib

import (
	"guardhouse/pkg/configkit"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type HostHeartbeatReq struct {
	IP          string `json:"ip" url:"ip"`
	Hostname    string `json:"hostname" url:"hostname"`
	Arch        string `json:"arch" url:"arch"`
	OS          string `json:"os" url:"os"`
	Kernel      string `json:"kernel" url:"kernel"`
	SystemInfo  string `json:"system_info" url:"system_info"`
	Environment string `json:"environment" url:"environment"`
}

type HostHeartbeatResp struct{}

func RandomHeartbeat() {
	for {
		t := rand.Intn(configkit.MaxSleepInterval)

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

	req := &HostHeartbeatReq{
		IP:          GetLocalIP(),
		Hostname:    GetHostname(),
		Arch:        GetArch(),
		OS:          GetOS(),
		Kernel:      GetKernel(),
		SystemInfo:  GetSystemInfo(),
		Environment: GetEnvironment(),
	}

	log.Printf("Heartbeat request: %+v", req)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(GetHostHeartbeatURL())
	if err != nil {
		log.Printf("Failed to send heartbeat: %v", err)

		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Response body: %s", resp.String())
		log.Printf("Heartbeat request failed with status: %d", resp.StatusCode())

		return
	}

	log.Printf("Heartbeat sent successfully")
}
