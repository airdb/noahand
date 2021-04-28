package noahlib

import (
	"airdb.io/airdb/sailor"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

type HostReq struct {
	IP string `url:"ip"`
	OS string `url:"os"`
	Hostname string `url:"hostname"`
	Timestamp string `url:"timestamp"`
	Arch string `url:"arch"`
	IsStart string `url:"is_start"`
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

func Heartbeat() {
	client := sailor.HTTPClient{}
	client.SetURL(GetConfigURL())
	// client.SetMethod(http.MethodGet)

	hostname, _ := os.Hostname()

	input := &HostReq{
		OS: runtime.GOOS,
		Arch: runtime.GOARCH,
		IP: GetLocalIP(),
		Hostname: hostname,
		Timestamp: fmt.Sprintf("", time.Now().Unix()),
	}

	client.SetBody(&input)

	client.SetUserAgent("noah-client/v0.0.1")

	var output HostResp

	err := sailor.HTTPRequest(&client, &output)
	if err != nil {
		log.Println(err)

		return
	}

	// log.Println("resp", output)
	log.Println("resp")
}
