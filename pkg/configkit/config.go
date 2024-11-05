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
