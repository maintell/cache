package cache

import (
	"strings"
	"time"
)

const (
	MetricTypeGetMem          = "get_mem"
	MetricTypeGetMemMiss      = "get_mem_miss"
	MetricTypeGetMemExpired   = "get_mem_expired"
	MetricTypeGetRedis        = "get_redis"
	MetricTypeGetRedisMiss    = "get_redis_miss"
	MetricTypeGetRedisExpired = "get_redis_expired"
	MetricTypeGetCache        = "get_cache"
	MetricTypeLoad            = "load"
	MetricTypeAsyncLoad       = "async_load"
	MetricTypeSetCache        = "set_cache"
	MetricTypeSetMem          = "set_mem"
	MetricTypeSetRedis        = "set_redis"
	MetricTypeDeleteMem       = "del_mem"
	MetricTypeDeleteRedis     = "del_redis"
)

type Metrics struct {
	// keys are namespacedKey, need trim namespace
	namespace string

	onMetric func(key string, metricType string, elapsedTime time.Duration)
}

func (m Metrics) Observe() func(string, interface{}, *error) {
	start := time.Now()
	return func(namespacedKey string, metricType interface{}, err *error) {
		if m.onMetric != nil {
			var metric string
			switch v := metricType.(type) {
			case *string:
				metric = *v
			case string:
				metric = v
			default:
				return
			}

			// ignore metric for error case
			if err != nil && *err != nil {
				return
			}
			key := strings.TrimPrefix(namespacedKey, m.namespace+":")
			m.onMetric(key, metric, time.Since(start))
		}
	}
}
