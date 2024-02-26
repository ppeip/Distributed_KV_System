package singleflight

import "sync"

// 代表正在进行中，或已经结束的请求。
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 管理不同key的请求call
type Group struct {
	mu sync.Mutex //protects m
	m  map[string]*call
}

// 对一个key的请求记录
// 请求的缓存机制，如果请求一个key发现前面有该请求了，就共享第一个请求的返回值，理解为合并同类请求
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 通过m里面是指针使得call的修改对所有的goroutine可见
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理

	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}
