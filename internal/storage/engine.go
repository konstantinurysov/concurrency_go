package storage

import (
	"sync"
)

type Engine struct {
	storage map[string]string
	mu      sync.RWMutex
}

func (e *Engine) Get(key string) string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.storage[key]
}

func (e *Engine) Set(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.storage[key] = value
}

func (e *Engine) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.storage, key)
}

func NewEngine() *Engine {
	return &Engine{
		storage: make(map[string]string),
		mu:      sync.RWMutex{},
	}
}
