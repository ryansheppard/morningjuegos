package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	Ctx    context.Context
}

func New(ctx context.Context, address string, database int) *Cache {
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
		Ctx:    ctx,
	}
}

func (c *Cache) SetKey(key string, value interface{}, ttl int) error {
	if c == nil {
		slog.Info("No cache provided, skipping SetKey", "key", key)
		return nil
	}

	seconds := time.Duration(ttl) * time.Second

	return c.Client.Set(c.Ctx, key, value, seconds).Err()
}

func (c *Cache) GetKey(key string) (interface{}, error) {
	if c == nil {
		slog.Info("No cache provided, not setting key in cache", "key", key)
		return nil, nil
	}

	result, err := c.Client.Get(c.Ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Cache) DeleteKey(key string) {
	if c == nil {
		slog.Info("No cache provided, not deleting key from cache", "key", key)
		return
	}

	c.Client.Del(c.Ctx, key)
}
