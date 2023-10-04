package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type Cache struct {
	Client *redis.Client
	Ctx    context.Context
	tracer trace.Tracer
}

func New(ctx context.Context, address string, database int, tracer trace.Tracer) *Cache {
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
		tracer: tracer,
	}
}

func (c *Cache) SetKey(ctx context.Context, key string, value interface{}, ttl int) error {
	if c == nil {
		slog.Info("No cache provided, skipping SetKey", "key", key)
		return nil
	}

	_, span := c.tracer.Start(ctx, "set-key")
	defer span.End()

	seconds := time.Duration(ttl) * time.Second

	return c.Client.Set(c.Ctx, key, value, seconds).Err()
}

func (c *Cache) GetKey(ctx context.Context, key string) (interface{}, error) {
	if c == nil {
		slog.Info("No cache provided, not setting key in cache", "key", key)
		return nil, nil
	}

	_, span := c.tracer.Start(ctx, "get-key")
	defer span.End()

	result, err := c.Client.Get(c.Ctx, key).Result()
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

	_, span := c.tracer.Start(ctx, "delete-key")
	defer span.End()

	c.Client.Del(c.Ctx, key)
}
