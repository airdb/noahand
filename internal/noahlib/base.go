package noahlib

import (
	"path"
	"runtime"
)

const (
	DefaultWorkDir = "/opt/noah"
	DarwinWorkDir  = "./noah-agent"
)

func GetWorkingPath() string {
	switch runtime.GOOS {
	case "linux":
		return DefaultWorkDir
	case "darwin":
		return DefaultWorkDir
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
