package configkit

import (
	"guardhouse/pkg/version"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	HomeDir            = "/opt/noah/"
	PluginDir          = "./plugins"
	ModulesDir         = "./modules"
	RunMode            = "local"
	DefaultDomain      = "https://aid.run"
	DefaultLocalDomain = "http://127.0.0.1:8080"

	// If agent counts less than 1k, the max sleep interval is 300s should be enough.
	// If agent count is more than 20K, suggest to set max sleep interval to 900s.
	MaxSleepInterval = 300

	// Linux
	SystemdFilepath = "/etc/systemd/system/noah.service"
	SystemdFilename = "noah.service"

	// MacOS
	LaunchctlFilepath = "/Library/LaunchDaemons/com.example.noah.plist"
	LaunchctlFilename = "com.example.noah.plist"

	TmpDir = "/tmp/noah/"

	AdminAddr = "127.0.0.1:403"
)

var AdminApiList = []string{
	"/heath",
	"/internal/noah/host",
	"/internal/noah/selfupdate",
	"/internal/noah/selfupgrade",
	"/internal/noah/download_plugin",
	"/internal/noah/cmd",
	"/internal/noah/exec",
}

var PluginMap = make(map[string]string)

func InitConfig() {
	if IsRunGoBuild() {
		HomeDir = "./"
		RunMode = "local"
		DefaultDomain = DefaultLocalDomain
		PluginDir = "./output/plugins"
		ModulesDir = "./output/modules"
	}

	// Developer mode
	if runtime.GOOS == "Darwin" {
		MaxSleepInterval = 30
	}

	GlobalConfig.RunMode = RunMode

	runtimeInfo := DefaultDomain
	GlobalConfig.Runtime = runtimeInfo

	GlobalConfig.BuildInfo = version.GetBuildInfo()

	GlobalConfig.AdminApiList = AdminApiList

	PluginMap = CheckPlugins(PluginDir)
}

// There is 2 ways to run a go program:
// 1. go run main.go
// 2. go build -o main main.go && ./main
// This function is used to check which way the program is running
// Because the diffent way will use different config path.
func IsRunGoBuild() bool {
	path := os.Args[0]

	if strings.Contains(path, "go-build") {
		log.Println("Program is running with go run")
	} else {
		log.Println("Program was built with go build")
	}

	return strings.Contains(path, "go-build")
}

// Check how many plugins in plugins directory.
func CheckPlugins(pluginDir string) map[string]string {
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".so") {
			pluginName := strings.TrimSuffix(info.Name(), ".so")
			tmps := strings.Split(pluginName, "_")
			pluginName = tmps[1] + tmps[0]

			titleCaser := cases.Title(language.Und)
			PluginMap[path] = titleCaser.String(pluginName)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for k, v := range PluginMap {
		println(k, v)
	}

	return PluginMap
}

func GetConfigURL() string {
	return path.Join(DefaultDomain, "/host")
}

func GetAdminHeathURL() string {
	return "http://" + AdminAddr + "/health"
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
