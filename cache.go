package main

import (
	"sync"
	"time"
)

type Cache struct {
	mutex   sync.Mutex
	items   map[string]*item
	ttl     int64
	cleaned int64
}

type item struct {
	value  int
	access int64
}

func NewCache(ttl int64) *Cache {
	c := Cache{items: make(map[string]*item), ttl: ttl, cleaned: time.Now().Unix()}
	return &c
}

func (c *Cache) Get(key string) (int, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().Unix()
	if c.cleaned < now-c.ttl {
		c.clean(now)
	}

	it, ok := c.items[key]
	if !ok {
		return 0, false
	}

	it.access = now
	return it.value, true
}

func (c *Cache) Set(key string, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().Unix()
	if c.cleaned < now-c.ttl {
		c.clean(now)
	}

	it, ok := c.items[key]
	if !ok {
		it = &item{}
		c.items[key] = it
	}

	it.value = value
	it.access = now
}

func (c *Cache) clean(now int64) {
	for key, it := range c.items {
		if it.access < now-c.ttl {
			delete(c.items, key)
		}
	}
	c.cleaned = now
}
