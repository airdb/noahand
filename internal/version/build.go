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

const nullVersion = "0.0.0"

// Build version info.
type BuildInfo struct {
	GoVersion string    `json:"go_version"`
	Env       string    `json:"env"`
	Repo      string    `json:"repo"`
	Version   string    `json:"version"`
	Build     string    `json:"build"`
	BuildTime string    `json:"build_time"`
	CreatedAt time.Time `json:"created_at"`
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

func getDeployVersion(executable string) string {
	var stringBuilder strings.Builder

	cmd := exec.CommandContext(context.Background(), executable, "-version")
	cmd.Stdout = &stringBuilder
	cmd.Stderr = &stringBuilder
	err := cmd.Run()
	if err != nil {
		return nullVersion
	}

	return stringBuilder.String()
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
