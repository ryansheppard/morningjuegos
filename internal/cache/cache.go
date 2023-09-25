package cache

import "github.com/redis/go-redis/v9"

var cache *redis.Client

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

func Get() *redis.Client {
	return cache
}
