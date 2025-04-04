package redis

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/gomodule/redigo/redis"
)

func ConnectRedis(cfg config.RedisConfig) (redis.Conn, error) {
	address := fmt.Sprintf("%s:%s", cfg.REDIS_HOST, cfg.REDIS_PORT)
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
