package main

import (
	"context"
	"sync"
	"time"
)

type KVDatabase interface {
	Get(key string) (string, error)
	Keys() ([]string, error)
	MGet(keys []string) ([]*string, error)
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func LongFunction() {
	// This function is supposed to be long-running
	// and should be optimized to be concurrent
}

func LongFunctionWrapper(ctx context.Context, funcToRun func()) {
	// This function should run the LongFunction
	// with a timeout of 5 seconds
	for {
		select {
		case <-ctx.Done():
			return
		default:
			funcToRun()
		}
	}
}

func main() {
	//run long function with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	LongFunctionWrapper(ctx, LongFunction)

}
