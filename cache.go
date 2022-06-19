package main

import (
	"sync"
)

type cache struct {
	mu    sync.Locker
	cache map[string]cacheEntity
}

func NewCache() cache {
	return cache{mu: &sync.Mutex{}, cache: make(map[string]cacheEntity)}
}

type cacheEntity struct {
	mu sync.Locker

	timeout float64

	val interface{}
}

//处理过期的key
func (c *cache) handleOverkey() {

}

func (c *cache) Put(key string, val interface{}, time float64) {
	node := cacheEntity{
		mu:      &sync.Mutex{},
		timeout: float64(time),
		val:     val,
	}
	c.mu.Lock()

	c.cache[key] = node

	c.mu.Unlock()

}

func (c *cache) Get(key string) interface{} {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.cache[key]

}
