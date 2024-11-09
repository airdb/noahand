package configkit

import (
	"log"

	"github.com/spf13/viper"
)

var CfgFile string

type Config struct {
	RunMode      string   `mapstructure:"run_mode" json:"run_mode" default:"local"`
	Runtime      string   `mapstructure:"runtime" json:"runtime" default:"runtime"`
	BuildInfo    any      `mapstructure:"build_info" json:"build_info"`
	AdminApiList []string `mapstructure:"admin_api_list" json:"admin_api_list"`
}

var GlobalConfig Config

const (
	PluginStatusEnabled  = "enabled"
	PluginStatusDisabled = "disabled"
	PluginStatusUnknown  = "unknown"
	PluginStatusFailed   = "failed"
	PluginStatusOutdated = "outdated"

	ModuleStatusEnabled  = "enabled"
	ModuleStatusDisabled = "disabled"
	ModuleStatusUnknown  = "unknown"
	ModuleStatusFailed   = "failed"
	ModuleStatusRunning  = "running"
	ModuleStatusOutdated = "outdated"
)

type RuntimeConfig struct {
	HomeDir   string `mapstructure:"home_dir" json:"home_dir"`
	GOOS      string `mapstructure:"goos" json:"goos"`
	GOARCH    string `mapstructure:"goarch" json:"goarch"`
	Version   string `mapstructure:"version" json:"version"`
	BuildInfo string `mapstructure:"build_info" json:"build_info"`
}

type MainConfig struct {
	RunMode string `mapstructure:"run_mode" json:"run_mode"`
}

type PluginConfig struct {
	PluginName     string `mapstructure:"plugin_name" json:"plugin_name"`
	PluginPath     string `mapstructure:"plugin_path" json:"plugin_path"`
	PluginVersion  string `mapstructure:"plugin_version" json:"plugin_version"`
	PluginPriority int    `mapstructure:"plugin_priority" json:"plugin_priority"`
	PluginStatus   string `mapstructure:"plugin_status" json:"plugin_status"`
}

type ModuleConfig struct {
	ModuleName     string `mapstructure:"module_name" json:"module_name"`
	ModulePath     string `mapstructure:"module_path" json:"module_path"`
	ModuleVersion  string `mapstructure:"module_version" json:"module_version"`
	ModulePriority int    `mapstructure:"module_priority" json:"module_priority"`
	ModuleStatus   string `mapstructure:"module_status" json:"module_status"`
}

func Init() {
	// log.Println("Load config ", CfgFile)
	viper.SetConfigName(CfgFile)

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(HomeDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Configfile not found: " + CfgFile)
		}
	}

	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Println("Using config file:", viper.ConfigFileUsed())
	log.Println("Read config from file: ", GlobalConfig.RunMode)
}

func GetConfig() *Config {
	return &GlobalConfig
}
