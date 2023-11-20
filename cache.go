package cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

// ISource
type ISource interface {
	Get(ctx context.Context, key string) (interface{}, error)
}

// ICache
type ICache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, data interface{}) error
}

// Cache
type Cache struct {
	cache  ICache
	source ISource
	sg     singleflight.Group
	ttl    time.Duration
}

func NewCache(opts ...CacheOption) *Cache {
	c := &Cache{}
	for _, opt := range opts {
		opt(c)
	}
	mc := &MemoryCache{ttl: c.ttl, slots: make(map[string]*MemoryCacheSlot)}
	mc.init()
	c.cache = mc
	return c
}

func (c *Cache) RGet(ctx context.Context, key string) *Result {
	data, err := c.Get(ctx, key)
	return &Result{data: data, err: err}
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	rst, err, _ := c.sg.Do(key, func() (interface{}, error) {
		data, err := c.cache.Get(ctx, key)
		switch err {
		case nil:
			return data, nil
		case errExpire, errNotFound:
			return c.get(ctx, key)
		default:
			return nil, errInternal
		}
	})
	return rst, err
}

func (c *Cache) get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.source.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if err := c.cache.Set(ctx, key, data); err != nil {
		return nil, err
	}
	return data, nil
}
