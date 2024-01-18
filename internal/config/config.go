package config

import (
	"sync"
	"testProject/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Port     int    `yaml:"port"`
	} `yaml:"db"`

	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance

}
