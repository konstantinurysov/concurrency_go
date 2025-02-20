package main

import "sync"

type KVDatabase interface {
	Get(key string) (string, error)
	Keys() ([]string, error)
	MGet(keys []string) ([]*string, error)
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}
