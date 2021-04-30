package noahlib

import (
	"airdb.io/airdb/sailor"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/user"
	"runtime"
	"time"
)

type HostReq struct {
	IP string `url:"ip"`
	OS string `url:"os"`
	Hostname string `url:"hostname"`
	Timestamp string `url:"timestamp"`
	Arch string `url:"arch"`
	IsStart string `url:"is_start,omitempty"`
	Username string `url:"username"`
}

type HostResp struct {
}

func GetConfigURL() string {
	configURL := "http://sg.airdb.host/host"

	return configURL
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
		maxSleepInterval := 60

		// nolint: gosec
		t := rand.Intn(maxSleepInterval)

		log.Println("sleep time seed", t)
		t =10

		time.Sleep(time.Duration(t) * time.Second)

		Heartbeat()
	}
}

func Heartbeat() {
	client := sailor.HTTPClient{}
	client.SetURL(GetConfigURL())
	// client.SetMethod(http.MethodGet)

	hostname, _ := os.Hostname()
	user, _ := user.Current()

	input := &HostReq{
		OS: runtime.GOOS,
		Arch: runtime.GOARCH,
		IP: GetLocalIP(),
		Hostname: hostname,
		Timestamp: fmt.Sprintf("%v", time.Now().Unix()),
		Username: user.Username,
	}

	client.SetBody(&input)

	client.SetUserAgent("noah-client/v0.0.1")

	var output HostResp

	err := client.HTTPRequest(&client, &output)
	if err != nil {
		log.Println(err)
	}

	// log.Println("resp", output)
	log.Println("resp")
}
