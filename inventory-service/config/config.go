package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort string `mapstructure:"GRPC_PORT"`
	GrpcHost string `mapstructure:"GRPC_HOST"`
	AppEnv   string `mapstructure:"APP_ENV"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
