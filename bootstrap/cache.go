package bootstrap

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisCache(config *Config) *redis.Client {
	r := config.Redis
	addr := fmt.Sprintf("%s:%s", r.Host, r.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: r.Password,
	})

	return rdb
}
