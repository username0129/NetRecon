package test

import (
	"backend/internal/core"
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	cache := core.InitializeCache()
	_ = cache.Set("key", []byte("hello"))
	data1, err := cache.Get("key")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("获取结果：%v\n", string(data1))

	_ = cache.Set("key", []byte("world"))
	data2, err := cache.Get("key")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("获取结果：%v\n", string(data2))
}
