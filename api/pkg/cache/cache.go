package cache

import (
	"sync"
	"time"

	"github.com/bxrne/beacon/api/internal/config"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items sync.Map
	mu    sync.RWMutex
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	expiration := time.Now().Add(ttl).UnixNano()
	c.items.Store(key, CacheItem{
		Value:      value,
		Expiration: expiration,
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, exists := c.items.Load(key)
	if !exists {
		return nil, false
	}

	cacheItem := item.(CacheItem)
	if time.Now().UnixNano() > cacheItem.Expiration {
		c.items.Delete(key)
		return nil, false
	}

	return cacheItem.Value, true
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

func (c *Cache) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			now := time.Now().UnixNano()
			c.items.Range(func(key, value interface{}) bool {
				item := value.(CacheItem)
				if now > item.Expiration {
					c.items.Delete(key)
				}
				return true
			})
		}
	}()
}

func New(cfg *config.Config) *Cache {
	cache := &Cache{}
	cache.StartCleanup(time.Duration(cfg.Server.CacheTTL) * time.Second)
	return cache
}
