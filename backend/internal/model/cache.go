package model

type CacheItem struct {
	Value      []byte `json:"value,omitempty"`      // 值
	Expiration int64  `json:"expiration,omitempty"` // 过期时间
}
