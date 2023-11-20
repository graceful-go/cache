package cache

import "time"

type MemoryCacheSlot struct {
	t time.Time
	d interface{}
}

func (m MemoryCacheSlot) IsExpire(ttl time.Duration) bool {
	return time.Since(m.t) > ttl
}

func (m MemoryCacheSlot) Data() interface{} {
	return m.d
}
