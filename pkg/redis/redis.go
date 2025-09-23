package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/t1xelLl/projectWithOrder/configs"
)

func NewClientRedis(cfg configs.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
