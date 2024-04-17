package main

import "sync"

type ValKey struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewValKey() *ValKey {
	return &ValKey{
		data: map[string][]byte{},
	}
}

func (vk *ValKey) Set(key string, val []byte) error {
	vk.mu.Lock()

	defer vk.mu.Unlock()

	vk.data[key] = val
	return nil
}

func (vk *ValKey) Get(key string) ([]byte, bool) {
	vk.mu.RLock()

	defer vk.mu.RUnlock()
	val, ok := vk.data[key]
	return val, ok
}
