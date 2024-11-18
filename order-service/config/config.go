package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerPort        string `env:"SERVER_PORT" envDefault:"8082"`
	ServerHost        string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	AppEnv            string `env:"APP_ENV" envDefault:"development"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"debug"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	fmt.Println(cfg)
	return cfg, nil
}
