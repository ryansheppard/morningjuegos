package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var cache *redis.Client

var ctx = context.Background()

func New(address string) {
	if address == "" {
		return
	}

	if cache == nil {
		cache = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "",
			DB:       0,
		})
	}
}

func SetKey(key string, value interface{}, ttl int) error {
	if cache == nil {
		return nil
	}

	return cache.Set(ctx, key, value, ttl).Err()
}

func GetKey(key string) (interface{}, error) {
	if cache == nil {
		return nil, nil
	}

	result, err := cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
