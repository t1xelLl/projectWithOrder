package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const cacheTTL = 24 * time.Hour

type Cache struct {
	*redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) getKey(uid string) string {
	return "order:" + uid
}

func (c *Cache) SetData(ctx context.Context, uid string, value []byte) error {
	key := c.getKey(uid)
	return c.Set(ctx, key, value, cacheTTL).Err()
}

func (c *Cache) GetData(ctx context.Context, uid string) ([]byte, error) {
	key := c.getKey(uid)
	res, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

func (c *Cache) DelData(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

func (c *Cache) GetAllKeys(ctx context.Context) ([]string, error) {
	pattern := c.getKey("*")
	keys, err := c.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
