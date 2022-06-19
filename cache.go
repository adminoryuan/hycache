package main

import (
	"sync"
	"time"
)

type cache struct {
	mu    sync.RWMutex
	cache map[string]cacheEntity
}

func NewCache() cache {

	ca := cache{mu: sync.RWMutex{}, cache: make(map[string]cacheEntity)}
	go ca.handleOverkey()

	return ca
}

type cacheEntity struct {
	mu sync.RWMutex

	timeout int64

	val interface{}
}

//处理过期的key
func (c *cache) handleOverkey() {
	for {

		c.mu.RLock()
		temp := c.cache
		c.mu.RUnlock()
		for k, v := range temp {

			if GetCurrentTime() < int64(v.timeout) {
				continue
			}

			c.mu.Lock()

			delete(c.cache, k)

			c.mu.Unlock()

		}
		time.Sleep(600 * time.Microsecond)
	}
}
func GetCurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}
func (c *cache) Put(key string, val interface{}, Time int64) {
	start := GetCurrentTime()

	node := cacheEntity{
		mu:      sync.RWMutex{},
		timeout: start + (Time),
		val:     val,
	}
	c.mu.Lock()

	c.cache[key] = node

	c.mu.Unlock()

}

func (c *cache) Get(key string) interface{} {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.cache[key].val

}
