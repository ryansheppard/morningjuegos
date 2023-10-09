package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func New(address string, database int) *Cache {
	if address == "" {
		slog.Info("No redis address provided, skipping caching")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       database,
	})

	slog.Info("Created redis connection", "address", address, "database", database)

	return &Cache{
		Client: client,
	}
}

func (c *Cache) SetKey(ctx context.Context, key string, value interface{}, ttl int) error {
	if c == nil {
		slog.Info("No cache provided, skipping SetKey", "key", key)
		return nil
	}

	seconds := time.Duration(ttl) * time.Second

	return c.Client.Set(ctx, key, value, seconds).Err()
}

func (c *Cache) GetKey(ctx context.Context, key string) (interface{}, error) {
	if c == nil {
		slog.Info("No cache provided, not setting key in cache", "key", key)
		return nil, nil
	}

	result, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Cache) DeleteKey(ctx context.Context, key string) {
	if c == nil {
		slog.Info("No cache provided, not deleting key from cache", "key", key)
		return
	}

	c.Client.Del(ctx, key)
}
