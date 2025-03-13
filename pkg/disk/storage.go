package disk

import (
	"concurrency_hw1/pkg/common"
	"concurrency_hw1/pkg/logger"
	"fmt"
	"os"
	"sync"
	"time"
)

type DiskStorage struct {
	file      *os.File
	log       *logger.Logger
	path      string
	mu        sync.RWMutex
	batchSize int
}

func NewDiskStorage(path string, batchSize string, log *logger.Logger) (*DiskStorage, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	size, err := common.ParseSize(batchSize)
	if err != nil {
		return nil, err
	}

	return &DiskStorage{
		path:      path,
		file:      f,
		mu:        sync.RWMutex{},
		batchSize: size,
		log:       log,
	}, nil
}

func (d *DiskStorage) Append(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	info, err := d.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}

	if info.Size()+int64(len(data)) > int64(d.batchSize) {
		err = d.file.Close()
		if err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		timestamp := time.Now().Unix()
		newPath := fmt.Sprintf("%s.%d", d.path, timestamp)
		d.file, err = os.OpenFile(newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create new file: %w", err)
		}
		d.log.Info("Created new file: %s", newPath)
	}
	if err != nil {
		return fmt.Errorf("failed to append data to file: %w", err)
	}

	if _, err := d.file.Write(data); err != nil {
		return err
	}

	return d.file.Sync()
}

func (d *DiskStorage) Close() error {
	return d.file.Close()
}
