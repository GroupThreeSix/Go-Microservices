package config

import "github.com/caarlos0/env/v11"

type Config struct {
	ServerPort           string `env:"SERVER_PORT" envDefault:"8081"`
	ServerHost           string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	InventoryServiceHost string `env:"INVENTORY_SERVICE_HOST" envDefault:"inventory-service"`
	InventoryServicePort string `env:"INVENTORY_SERVICE_PORT" envDefault:"50051"`
	AppEnv               string `env:"APP_ENV" envDefault:"development"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
