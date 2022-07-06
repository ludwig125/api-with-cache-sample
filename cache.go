package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	item       bool
	lastAccess int64
}

type TTLMap struct {
	m map[string]*CacheItem
	l sync.Mutex
}

// func NewCache(ln int, maxTTL int) (m *TTLMap) {
func NewCache(maxTTL int) (m *TTLMap) {
	// m = &TTLMap{m: make(map[string]*CacheItem, ln)}
	m = &TTLMap{m: make(map[string]*CacheItem)}
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
	return
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k string, v bool) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &CacheItem{item: v}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap) Get(k string) (v bool) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.item
		it.lastAccess = time.Now().Unix()
	}
	m.l.Unlock()
	return

}
