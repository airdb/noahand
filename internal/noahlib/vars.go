package noahlib

import (
	"net"
	"net/url"
	"path"
	"runtime"
)

var (
	HomeDir       = "/opt/noah/"
	DefaultDomain = "https://aid.run"
	// DefaultDomain = "http://127.0.0.1:8080"
)

func GetConfigURL() string {
	return path.Join(DefaultDomain, "/host")
}

func GetHostHeartbeatURL() string {
	apiPath := "/apis/infra/v1/host/heartbeat"

	base, err := url.Parse(DefaultDomain)
	if err != nil {
		panic(err)
	}

	relative, err := url.Parse(apiPath)
	if err != nil {
		panic(err)
	}

	fullURL := base.ResolveReference(relative)

	return fullURL.String()
}

func GetBinPath() string {
	return path.Join(HomeDir, "bin")
}

func GetPluginPath() string {
	switch runtime.GOOS {
	case "linux":
		return path.Join(HomeDir, "plugins")
	case "darwin":
		return path.Join(HomeDir, "plugins")
	default:
		panic("unsupported os: " + runtime.GOOS)
	}
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

func GetHostname() string {
	return "hostname"
}

func GetArch() string {
	return runtime.GOARCH
}

func GetOS() string {
	return runtime.GOOS
}

func GetKernel() string {
	return runtime.GOOS
}

func GetSystemInfo() string {
	return "system info"
}

func GetEnvironment() string {
	return "env"
}
