package redis

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/gomodule/redigo/redis"
)

func ConnectRedis(cfg config.RedisConfig) (redis.Conn, error) {
	address := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redisConn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	_, err = redisConn.Do("PING")
	if err != nil {
		return nil, err
	}
	return redisConn, nil
}

func NewRedisPool(cfg config.RedisConfig) *redis.Pool {
	address := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
