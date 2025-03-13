package wal

import (
	"concurrency_hw1/pkg/logger"
	"context"
	"time"
)

type WALService struct {
	wal        WAL
	logger     *logger.Logger
	size       int
	timeout    time.Duration
	WALChannel chan ([]byte)
	batch      [][]byte
}

func NewWALService(wal WAL, size int, timeout time.Duration, logger *logger.Logger) *WALService {
	return &WALService{
		wal:        wal,
		size:       size,
		timeout:    timeout,
		logger:     logger,
		WALChannel: make(chan []byte),
		batch:      make([][]byte, 0),
	}
}

func (w *WALService) Start(ctx context.Context) {
	go w.run(ctx)
}

func (w *WALService) run(ctx context.Context) {
	t := time.NewTicker(w.timeout)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case operation := <-w.WALChannel:
			if len(w.batch) > w.size {
				w.handleRecord(w.batch)
				w.batch = nil
			}
			w.batch = append(w.batch, operation)
		case <-t.C:
			if len(w.batch) > 0 {
				w.handleRecord(w.batch)
				w.batch = nil
			}
		}
	}
}

// handleRecord writes the operations to the WAL
func (w *WALService) handleRecord(operations [][]byte) {
	for _, o := range operations {
		if err := w.wal.Append(o); err != nil {
			w.logger.Error("failed to record to wal: %w", err)
		}
	}
}
