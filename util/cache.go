package util

import (
	"GoProject/global"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
)

var cacheAdapter *cache.Cache

func InitCache() {
	global.Logger.Debug("init cache")
	cacheAdapter = cache.New(global.DefaultCacheExpiredTime, time.Minute*10)
}

func SetCache(key string, value interface{}, d time.Duration) {
	cacheAdapter.Set(key, value, d)
	global.Logger.WithFields(logrus.Fields{
		"key":   key,
		"value": value,
	}).Debug("set cache")
}

func SetCacheWithDefault(key string, value interface{}) {
	SetCache(key, value, global.DefaultCacheExpiredTime)
}

func GetCache(key string) (value interface{}, exists bool) {
	global.Logger.WithFields(logrus.Fields{
		"key": key,
	}).Debug("get cache")
	return cacheAdapter.Get(key)
}

func RemoveCache(key string) {
	cacheAdapter.Delete(key)
}
