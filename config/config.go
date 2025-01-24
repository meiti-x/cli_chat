package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Logger struct {
	Level string
	Path  string
}

type Server struct {
	Host string
	Port int
}

type Nats struct {
	ConnString string
}

type Database struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

type Redis struct {
	Host     string
	Password string
	DB       int
}

type Config struct {
	Server
	Nats
	Redis
	Logger
	Database
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
