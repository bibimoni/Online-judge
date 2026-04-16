package database

import (
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() (*redis.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Uri,
		Password: cfg.Redis.Password,
	})

	return rdb, nil
}
