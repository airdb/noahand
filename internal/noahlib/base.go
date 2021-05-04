package noahlib

import "runtime"

func GetPluginPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/srv/noah/plugins/"
	case "darwin":
		return "./plugins/"
	}

	return ""
}

const (
	NoahPath       = "/srv/noah"
	NoahBinPath    = "/srv/noah/bin"
	NoahModulePath = "/srv/noah/module/"
	NoahPluginPath = "/srv/noah/plugin/"
)
