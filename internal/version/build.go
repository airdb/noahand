package version

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Build version info.
type BuildInfo struct {
	GoVersion string
	Env       string
	Repo      string
	Version   string
	Build     string
	BuildTime string
	CreatedAt time.Time
}

var (
	Repo      string
	Version   string
	Build     string
	BuildTime string
	CreatedAt time.Time
)

func GetBuildInfo() *BuildInfo {
	return &BuildInfo{
		GoVersion: runtime.Version(),
		Env:       os.Getenv("ENV"),
		Repo:      Repo,
		Version:   Version,
		Build:     Build,
		BuildTime: BuildTime,
		CreatedAt: CreatedAt,
	}
}

func (info *BuildInfo) ToString() string {
	out, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return string(out)
}

func ToString() string {
	info := GetBuildInfo()
	out, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return string(out)
}

const nullVersion  = "0.0.0"
func getDeployVersion(executable string) string {
	var sb strings.Builder
	cmd := exec.CommandContext(context.Background(), executable, "-version")
	cmd.Stdout = &sb
	cmd.Stderr = &sb
	err := cmd.Run()
	if err != nil {
		return nullVersion
	}

	return sb.String()
}

func GetDeployVersion() string {
	executable, err := os.Executable()
	if err != nil {
		return nullVersion
	}

	return getDeployVersion(executable)
}

func GetRunningVersion() string {
	return Version
}
