package main

import (
	"context"
	"sync"
	"time"
)

/*
参考
https://stackoverflow.com/questions/25484122/map-with-ttl-option-in-go
https://bmf-tech.com/posts/Golang%E3%81%A7%E3%82%A4%E3%83%B3%E3%83%A1%E3%83%A2%E3%83%AA%E3%81%AA%E3%82%AD%E3%83%A3%E3%83%83%E3%82%B7%E3%83%A5%E3%82%92%E5%AE%9F%E8%A3%85%E3%81%99%E3%82%8B
https://hackernoon.com/in-memory-caching-in-golang

*/

type CacheItem struct {
	value      bool
	lastAccess int64
}

type TTLMapCache struct {
	m map[ID]*CacheItem
	l sync.Mutex
}

// func NewCache(ln int, maxTTL int) (m *TTLMap) {
func NewCache(maxTTL int) CacheRepository {
	m := &TTLMapCache{m: make(map[ID]*CacheItem)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return m
}

func (c *TTLMapCache) Put(ctx context.Context, k ID, v bool) error {
	c.l.Lock()
	defer c.l.Unlock()
	it, ok := c.m[k]
	if !ok {
		it = &CacheItem{value: v}
		c.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	return nil
}

func (c *TTLMapCache) Get(ctx context.Context, k ID) (bool, error) {
	c.l.Lock()
	defer c.l.Unlock()
	// if it, ok := m.m[k]; ok {
	if _, ok := c.m[k]; ok {
		// v := it.value
		return true, nil
		// it.lastAccess = time.Now().Unix() // GetでもlastAccessを更新してしまうとGetリクエストが来る限りキャッシュがExpireしない
	}
	return false, nil
}
