package noahlib

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
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
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error:", err)
		return "unknown"
	}
	fmt.Println("Hostname:", hostname)

	return hostname
}

func GetArch() string {
	return runtime.GOARCH
}

func GetOS() string {
	var info string

	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("bash", "-c", "grep '^PRETTY_NAME=' /etc/os-release | cut -d '\"' -f 2")
		output, err := cmd.Output()
		if err != nil {
			log.Println("Error:", err)
			info = "unknown"
		} else {
			info = strings.TrimSpace(string(output))
		}
	case "darwin":
		cmd := exec.Command("sw_vers")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Println("Error:", err)
			info = "unknown"
		} else {
			lines := strings.Split(out.String(), "\n")
			var productName, productVersion, buildVersion string
			for _, line := range lines {
				if strings.HasPrefix(line, "ProductName:") {
					productName = strings.TrimSpace(strings.Split(line, ":")[1])
				} else if strings.HasPrefix(line, "ProductVersion:") {
					productVersion = strings.TrimSpace(strings.Split(line, ":")[1])
				} else if strings.HasPrefix(line, "BuildVersion:") {
					buildVersion = strings.TrimSpace(strings.Split(line, ":")[1])
				}
			}

			info = productName + " " + productVersion + " " + buildVersion
		}

	default:
		info = "unknown os: " + runtime.GOOS
	}

	return info
}

func GetKernel() string {
	return runtime.GOOS
}

func GetSystemInfo() string {
	var info string

	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("uname", "-a")
		output, err := cmd.Output()
		if err != nil {
			log.Println("Error:", err)
			info = "unknown"
		} else {
			info = string(output)
		}
	case "darwin":
		cmd := exec.Command("sw_vers")
		output, err := cmd.Output()
		if err != nil {
			log.Println("Error:", err)
			info = "unknown"
		} else {
			info = string(output)
		}
	default:
		log.Println("unsupported os: " + runtime.GOOS)
		info = "unknown os: " + runtime.GOOS
	}

	return info
}

func GetEnvironment() string {
	return strings.Join(os.Environ(), "\n")
}
