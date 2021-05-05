package noahlib

import (
	"path"
	"runtime"
)

const (
	LinuxWorkDir  = "/srv/noah"
	DarwinWorkDir = "./noah-agent"
)

func GetWorkingPath() string {
	switch runtime.GOOS {
	case "linux":
		return LinuxWorkDir
	case "darwin":
		return DarwinWorkDir
	default:
		return DarwinWorkDir
	}
}

func GetNoahBinPath() string {
	return path.Join(GetWorkingPath(), "module")
}

func GetNoahModulePath() string {
	return path.Join(GetWorkingPath(), "module")
}

func GetNoahPluginPath() string {
	return path.Join(GetWorkingPath(), "plugin")
}
