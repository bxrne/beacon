package metrics

import (
	"time"

	"github.com/bxrne/beacon/api/pkg/cache"
	"github.com/bxrne/beacon/api/pkg/config"
)

type MetricsCache struct {
	cache *cache.Cache
	ttl   time.Duration
}

func NewMetricsCache(cfg *config.Config) *MetricsCache {
	return &MetricsCache{
		cache: cache.New(cfg),
		ttl:   time.Duration(cfg.Server.CacheTTL) * time.Second,
	}
}

func (mc *MetricsCache) SetMetrics(hostname string, metrics DeviceMetrics) {
	mc.cache.Set(hostname, metrics, mc.ttl)
}

func (mc *MetricsCache) GetMetrics(hostname string) (DeviceMetrics, bool) {
	value, exists := mc.cache.Get(hostname)
	if !exists {
		return DeviceMetrics{}, false
	}
	return value.(DeviceMetrics), true
}

func (mc *MetricsCache) DeleteMetrics(hostname string) {
	mc.cache.Delete(hostname)
}
