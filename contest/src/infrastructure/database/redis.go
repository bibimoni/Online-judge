package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient(host, port, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	return rdb, nil
}
