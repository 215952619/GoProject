package util

import (
	"GoProject/global"
	"github.com/patrickmn/go-cache"
	"time"
)

var cacheAdapter *cache.Cache

func InitCache() {
	global.Logger.Debug("init cache")
	cacheAdapter = cache.New(global.DefaultCacheExpiredTime, time.Minute*10)
}

func SetCache(key string, value interface{}, d time.Duration) {
	cacheAdapter.Set(key, value, d)
}

func SetCacheWithDefault(key string, value interface{}) {
	cacheAdapter.Set(key, value, global.DefaultCacheExpiredTime)
}

func GetCache(key string) (interface{}, bool) {
	return cacheAdapter.Get(key)
}

func RemoveCache(key string) {
	cacheAdapter.Delete(key)
}
