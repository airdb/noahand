package noahlib

import (
	"path"
	"runtime"
)

var (
	HomeDir       = "/opt/noah/"
	DefaultDomain = "http://sg.airdb.host:8000"
)

func GetConfigURL() string {
	return DefaultDomain + "/host"
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
		return ""
	}
}
