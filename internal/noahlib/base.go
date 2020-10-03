package noahlib

import "runtime"

func GetPluginPath() string{
	switch runtime.GOOS {
	case "linux":
		return "/srv/noah/plugins/"
	case "darwin":
		return "./plugins/"
	}

	return ""
}
