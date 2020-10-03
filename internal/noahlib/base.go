package noahlib

import "runtime"

func GetPluginPath() string{
	switch runtime.GOOS {
	case "linux":
		return "/noah/plugins/"
	case "darwin":
		return "./plugins/"
	}

	return ""
}
