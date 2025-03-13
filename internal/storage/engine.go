package storage

import (
	"sync"
)

type EngineInterface interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
}

type Engine struct {
	storage map[string]string
	mu      sync.RWMutex
}

func (e *Engine) Get(key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	value, ok := e.storage[key]
	return value, ok
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
