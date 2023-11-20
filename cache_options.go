package cache

import "time"

type CacheOption func(*Cache)

func WithTTL(ttl time.Duration) func(*Cache) {
	return func(c *Cache) {
		c.ttl = ttl
	}
}

func WithSource(source ISource) func(*Cache) {
	return func(c *Cache) {
		c.source = source
	}
}
