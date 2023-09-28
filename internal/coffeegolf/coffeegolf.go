package coffeegolf

import (
	"github.com/ryansheppard/morningjuegos/internal/cache"
)

type CoffeeGolf struct {
	Query *Query
	Cache *cache.Cache
}

func NewCoffeeGolf(query *Query, cache *cache.Cache) *CoffeeGolf {
	return &CoffeeGolf{
		Query: query,
		Cache: cache,
	}
}
