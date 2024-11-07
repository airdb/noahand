package configkit

import (
	"guardhouse/pkg/oskit"
	"guardhouse/pkg/version"
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
	RunMode            = "local"
	DefaultDomain      = "https://aid.run"
	DefaultLocalDomain = "http://127.0.0.1:8080"
)

var AdminApiList = []string{
	"/internal/noah/host",
	"/internal/noah/selfupdate",
	"/internal/noah/selfupgrade",
	"/internal/noah/download_plugin",
	"/internal/noah/cmd",
	"/internal/noah/exec",
}

var PluginMap = make(map[string]string)

func InitConfig() {
	if oskit.IsRunGoBuild() {
		HomeDir = "./"
		RunMode = "local"
		DefaultDomain = DefaultLocalDomain
		PluginDir = "./output/plugins"
	}

	GlobalConfig.RunMode = RunMode

	runtimeInfo := DefaultDomain
	GlobalConfig.Runtime = runtimeInfo

	GlobalConfig.BuildInfo = version.GetBuildInfo()

	GlobalConfig.AdminApiList = AdminApiList

	PluginMap = CheckPlugins(PluginDir)
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
