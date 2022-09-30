package config

import (
	"fmt"
	"github.com/comoyi/valheim-syncer-server/log"
	"github.com/comoyi/valheim-syncer-server/util/fsutil"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
)

var Conf Config

type Config struct {
	LogLevel string `toml:"log_level" mapstructure:"log_level"`
	Gui      string `toml:"gui" mapstructure:"gui"`
	Port     int    `toml:"port" mapstructure:"port"`
	Dir      string `toml:"dir" mapstructure:"dir"`
	Interval int64  `toml:"interval" mapstructure:"interval"`
}

func initDefaultConfig() {
	viper.SetDefault("log_level", log.Off)
	viper.SetDefault("gui", "ON")
	viper.SetDefault("port", 8080)
	viper.SetDefault("dir", "")
	viper.SetDefault("interval", 10)
}

func LoadConfig() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("%s%s%s", "$HOME", string(os.PathSeparator), ".valheim-syncer-server"))

	initDefaultConfig()

	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("Read config failed, err: %v\n", err)
		//return
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Errorf("Unmarshal config failed, err: %v\n", err)
		return
	}
	log.Debugf("config: %+v\n", Conf)
}

var saveMutex = &sync.Mutex{}

func SaveConfig() error {
	saveMutex.Lock()
	defer saveMutex.Unlock()

	err := viper.WriteConfig()
	if err == nil {
		return nil
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("Get os.UserHomeDir failed, err: %v\n", err)
		return err
	}
	log.Debugf("userHomeDir: %s\n", userHomeDir)

	configPath := filepath.Join(userHomeDir, ".valheim-syncer-server")
	configFile := filepath.Join(configPath, "config.toml")
	log.Debugf("configFile: %s\n", configFile)

	exist, err := fsutil.Exists(configPath)
	if err != nil {
		log.Warnf("Check isPathExist failed, err: %v\n", err)
		return err
	}
	if !exist {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			log.Warnf("Get os.MkdirAll failed, err: %v\n", err)
			return err
		}
	}

	err = viper.WriteConfigAs(configFile)
	if err != nil {
		log.Errorf("WriteConfigAs failed, err: %v\n", err)
		return err
	}
	return nil
}
