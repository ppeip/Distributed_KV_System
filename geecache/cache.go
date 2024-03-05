package geecache

import (
	"geecache/lru"
	"sync"
)

// cache负责对lru的并发控制
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

// 并发控制

// 添加缓存
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	// 函数结束解锁
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// 通过key获取缓存值
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
