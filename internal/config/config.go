package config

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"github.com/spf13/viper"
)

type Config struct {
	Cors middleware.Cors
	Port string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
