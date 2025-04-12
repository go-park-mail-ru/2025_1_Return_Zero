package config

import (
	"os"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/spf13/viper"
)

type PostgresConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	MaxOpenConns      int `mapstructure:"max_open_conns"`
	MaxIdleConns      int `mapstructure:"max_idle_conns"`
	MaxLifetime       int `mapstructure:"max_lifetime"`
}

type S3Config struct {
	S3Region        string
	S3Endpoint      string
	S3TracksBucket  string
	S3ImagesBucket  string
	S3AccessKey     string
	S3SecretKey     string
	S3Duration      time.Duration `mapstructure:"s3_duration"`
}

type RedisConfig struct {
	RedisHost string
	RedisPort string
}

type Config struct {
	Cors       middleware.Cors
	Port       string `mapstructure:"port"`
	Pagination deliveryModel.PaginationConfig
	Postgres   PostgresConfig
	S3         S3Config
	Redis      RedisConfig
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

	config.Postgres.PostgresHost = os.Getenv("POSTGRES_HOST")
	config.Postgres.PostgresPort = os.Getenv("POSTGRES_PORT")
	config.Postgres.PostgresUser = os.Getenv("POSTGRES_USER")
	config.Postgres.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	config.Postgres.PostgresDB = os.Getenv("POSTGRES_DB")

	config.S3.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	config.S3.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	config.S3.S3Region = os.Getenv("S3_REGION")
	config.S3.S3Endpoint = os.Getenv("S3_ENDPOINT")
	config.S3.S3TracksBucket = os.Getenv("S3_TRACKS_BUCKET")
	config.S3.S3ImagesBucket = os.Getenv("S3_IMAGES_BUCKET")

	config.Redis.RedisHost = os.Getenv("REDIS_HOST")
	config.Redis.RedisPort = os.Getenv("REDIS_PORT")

	return &config, nil
}
