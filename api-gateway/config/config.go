package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerPort        string `env:"SERVER_PORT" envDefault:"8080"`
	ServerHost        string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	ProductServiceURL string `env:"PRODUCT_SERVICE_URL" envDefault:"http://product-service:8081"`
	OrderServiceURL   string `env:"ORDER_SERVICE_URL" envDefault:"http://order-service:8082"`
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
