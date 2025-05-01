package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type CSRFConfig struct {
	CSRFHeaderName  string `mapstructure:"csrf_header_name"`
	CSRFCookieName  string `mapstructure:"csrf_cookie_name"`
	CSRFTokenLength int    `mapstructure:"csrf_token_length"`
}

type PostgresConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	MaxOpenConns     int `mapstructure:"max_open_conns"`
	MaxIdleConns     int `mapstructure:"max_idle_conns"`
	MaxLifetime      int `mapstructure:"max_lifetime"`
}

type S3Config struct {
	S3Region       string
	S3Endpoint     string
	S3TracksBucket string
	S3ImagesBucket string
	S3AccessKey    string
	S3SecretKey    string
	S3Duration     time.Duration `mapstructure:"s3_duration"`
}

type RedisConfig struct {
	RedisHost string
	RedisPort string
}

type Cors struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type ArtistService struct {
	Port int `mapstructure:"port"`
	Host string
}

type AlbumService struct {
	Port int `mapstructure:"port"`
	Host string
}

type TrackService struct {
	Port int `mapstructure:"port"`
	Host string
}

type AuthService struct {
	Port int `mapstructure:"port"`
	Host string
}

type UserService struct {
	Port int `mapstructure:"port"`
	Host string
}

type Services struct {
	ArtistService ArtistService `mapstructure:"artist_service"`
	AlbumService  AlbumService  `mapstructure:"album_service"`
	TrackService  TrackService  `mapstructure:"track_service"`
	AuthService   AuthService   `mapstructure:"auth_service"`
	UserService   UserService   `mapstructure:"user_service"`
}

type PaginationConfig struct {
	MaxOffset     int `mapstructure:"max_offset"`
	MaxLimit      int `mapstructure:"max_limit"`
	DefaultOffset int `mapstructure:"default_offset"`
	DefaultLimit  int `mapstructure:"default_limit"`
}

type Config struct {
	Cors       Cors
	Port       int `mapstructure:"port"`
	Pagination PaginationConfig
	Postgres   PostgresConfig
	S3         S3Config
	Redis      RedisConfig
	CSRF       CSRFConfig
	Services   Services
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

	config.Services.ArtistService.Host = os.Getenv("ARTIST_SERVICE_HOST")
	config.Services.AlbumService.Host = os.Getenv("ALBUM_SERVICE_HOST")
	config.Services.TrackService.Host = os.Getenv("TRACK_SERVICE_HOST")
	config.Services.AuthService.Host = os.Getenv("AUTH_SERVICE_HOST")
	config.Services.UserService.Host = os.Getenv("USER_SERVICE_HOST")

	return &config, nil
}
