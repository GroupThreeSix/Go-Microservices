package config

import "github.com/caarlos0/env/v10"

type Config struct {
	GrpcPort string `env:"GRPC_PORT" envDefault:"50051"`
	GrpcHost string `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	AppEnv   string `env:"APP_ENV" envDefault:"development"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
