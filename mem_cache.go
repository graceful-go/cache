package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	errExpire   = errors.New("key expired")
	errNotFound = errors.New("key not found")
	errInternal = errors.New("cache internal error")
)

type MemoryCache struct {
	ttl    time.Duration
	slots  map[string]*MemoryCacheSlot
	slotsP sync.Pool
	mu     sync.RWMutex
}

func (m *MemoryCache) init() {
	m.slotsP = sync.Pool{
		New: func() any {
			return &MemoryCacheSlot{}
		},
	}
}

func (m *MemoryCache) Clear(ctx context.Context, key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	slot, ok := m.slots[key]
	if !ok {
		return
	}

	delete(m.slots, key)
	m.slotsP.Put(&slot)
}

func (m *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	slot, ok := m.slots[key]
	if !ok {
		return nil, errNotFound
	}

	if slot.IsExpire(m.ttl) {
		go m.Clear(context.Background(), key)
		return nil, errExpire
	}

	return slot.Data(), nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, data interface{}) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	slot := m.slotsP.Get().(*MemoryCacheSlot)
	slot.d = data
	slot.t = time.Now()

	m.slots[key] = slot

	return nil
}
