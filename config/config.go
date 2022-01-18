package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

//stores all conf values read by viper from config env files

func setupViperConfig() {
	configPath := os.Getenv("APP_DIR") + "/config"
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config.local")
}

func setRootIfNotExist() {
	_, present := os.LookupEnv("APP_DIR")
	if !present {
		_, b, _, _ := runtime.Caller(0)
		os.Setenv("APP_DIR", filepath.Dir(filepath.Dir(b)))
	}
}

func init() {
	setRootIfNotExist()
	setupViperConfig()
}

type AppConfig struct {
	PubSub       *PubSubCfg
	DeployConfig *DeployConfig
}

type DeployConfig struct {
	Port         string `mapstructure:"PORT"`
	AppDir       string `mapstructure:"APP_DIR"`
	NumGoRoutine int    `mapstructure:"NUM_GO_ROUTINE"`
}

type PubSubCfg struct {
	ProjectID string `mapstructure:"PUBSUB_PROJECT_ID"`
	SubID     string `mapstructure:"PUBSUB_SUB_ID"`
}

func LoadConfig() (config AppConfig, err error) {
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	var pubsub PubSubCfg
	var deploy DeployConfig

	err = viper.Unmarshal(&pubsub)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&deploy)
	if err != nil {
		return
	}

	config.PubSub = &pubsub
	config.DeployConfig = &deploy
	return
}
