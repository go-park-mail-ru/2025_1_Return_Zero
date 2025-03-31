package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	MaxOpenConns      int `mapstructure:"max_open_conns"`
	MaxIdleConns      int `mapstructure:"max_idle_conns"`
	MaxLifetime       int `mapstructure:"max_lifetime"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Config struct {
	Port     string `mapstructure:"port"`
	Postgres PostgresConfig
	Redis    RedisConfig
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	config.Postgres.POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	config.Postgres.POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	config.Postgres.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	config.Postgres.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	config.Postgres.POSTGRES_DB = os.Getenv("POSTGRES_DB")
	
	config.Redis.Host = os.Getenv("REDIS_HOST")
	config.Redis.Port = os.Getenv("REDIS_PORT")
	fmt.Printf("Config successfully loaded")
	return &config, nil
}
