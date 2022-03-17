package noahlib

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/airdb/sailor"
	"github.com/airdb/sailor/version"
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

type HostResp struct {
}

const (
	DefaultDomain = "http://sg.airdb.host"
)

func GetConfigURL() string {
	return DefaultDomain + "/host"
}

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
		maxSleepInterval := 300

		// nolint: gosec
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
	client := sailor.HTTPClient{}
	client.SetURL(GetConfigURL())

	hostname, _ := os.Hostname()
	user, _ := user.Current()

	input := &HostReq{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		IP:        GetLocalIP(),
		Hostname:  hostname,
		Timestamp: fmt.Sprintf("%v", time.Now().Unix()),
		Username:  user.Username,
		Version:   version.Version,
	}

	client.SetBody(&input)

	client.SetUserAgent(fmt.Sprintf("%s/%s", version.Repo, version.Version))

	var output HostResp

	err := client.HTTPRequest(&client, &output)
	if err != nil {
		log.Println(err)
	}
}
