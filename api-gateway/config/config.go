package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort        string `mapstructure:"SERVER_PORT"`
	ServerHost        string `mapstructure:"SERVER_HOST"`
	ProductServiceURL string `mapstructure:"PRODUCT_SERVICE_URL"`
	OrderServiceURL   string `mapstructure:"ORDER_SERVICE_URL"`
	AppEnv            string `mapstructure:"APP_ENV"`
	LogLevel          string `mapstructure:"LOG_LEVEL"`
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
