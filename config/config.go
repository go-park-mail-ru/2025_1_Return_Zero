package config

import (
	"fmt"
	"os"
	"time"
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

type S3Config struct {
	S3_REGION        string
	S3_ENDPOINT      string
	S3_TRACKS_BUCKET string
	S3_IMAGES_BUCKET string
	S3_ACCESS_KEY    string
	S3_SECRET_KEY    string
	S3_DURATION      time.Duration `mapstructure:"s3_duration"`
}	

type Config struct {
	Port     string `mapstructure:"port"`
	Postgres PostgresConfig
	Redis    RedisConfig
	S3       S3Config
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

	config.S3.S3_ACCESS_KEY = os.Getenv("S3_ACCESS_KEY")
	config.S3.S3_SECRET_KEY = os.Getenv("S3_SECRET_KEY")
	config.S3.S3_REGION = os.Getenv("S3_REGION")
	config.S3.S3_ENDPOINT = os.Getenv("S3_ENDPOINT")
	config.S3.S3_TRACKS_BUCKET = os.Getenv("S3_TRACKS_BUCKET")
	config.S3.S3_IMAGES_BUCKET = os.Getenv("S3_IMAGES_BUCKET")
	
	config.Redis.Host = os.Getenv("REDIS_HOST")
	config.Redis.Port = os.Getenv("REDIS_PORT")
	fmt.Println("Config successfully loaded")
	return &config, nil
}
