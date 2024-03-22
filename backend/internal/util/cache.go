package util

import (
	"backend/internal/e"
	"backend/internal/global"
	"backend/internal/model"
	"encoding/json"
	"strconv"
	"time"
)

func SetCacheItem(key string, data []byte, e time.Duration) {
	item := model.CacheItem{
		Value:      data,
		Expiration: time.Now().Add(e).Unix(),
	}
	itemBytes, _ := json.Marshal(item)
	_ = global.Cache.Set(key, itemBytes)
}

func GetCacheItem(key string) (model.CacheItem, error) {
	itemBytes, err := global.Cache.Get(key)
	if err != nil {
		return model.CacheItem{}, err
	}

	var item model.CacheItem
	_ = json.Unmarshal(itemBytes, &item)
	if time.Now().Unix() > item.Expiration {
		return model.CacheItem{}, e.ErrCacheEntryTimeout
	} else {
		return item, nil
	}
}

func ItemToInt(item model.CacheItem) int {
	if n, err := strconv.Atoi(string(item.Value)); err != nil {
		return 0
	} else {
		return n
	}
}
