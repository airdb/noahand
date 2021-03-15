package noahlib

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"airdb.io/airdb/sailor"
)

type HostReq struct {
	IP string
	OS string
}

type HostResp struct {
}

func GetConfigURL() string {
	configURL := "http://hk.airdb.host"

	hostname, _ := os.Hostname()

	goos := runtime.GOOS

	return fmt.Sprintf("%s/host?timestamp=%v&os=%v&ip=%v&hostname=%v",
		configURL,
		time.Now().Unix(),
		goos,
		GetLocalIP(),
		hostname,
	)
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
	client.SetMethod(http.MethodGet)

	params := url.Values{}
	client.SetValues(params)

	client.SetHeaders(map[string]string{
		"User-Agent": "noah-client",
	})

	var output HostResp

	err := sailor.HTTPRequest(&client, &output)
	if err != nil {
		log.Println(err)

		return
	}

	// log.Println("resp", output)
	log.Println("resp")
}
