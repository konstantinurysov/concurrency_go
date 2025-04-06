package disk

import (
	"concurrency_hw1/pkg/common"
	"concurrency_hw1/pkg/logger"
	"context"
	"fmt"
	"os"
	"time"
)

type DiskStorage struct {
	file      *os.File
	log       *logger.Logger
	path      string
	batchSize int
	storageCh chan []byte
}

func NewDiskStorage(path string, batchSize string, log *logger.Logger) (*DiskStorage, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	size, err := common.ParseSize(batchSize)
	if err != nil {
		return nil, err
	}

	return &DiskStorage{
		file:      file,
		log:       log,
		path:      path,
		batchSize: size,
		storageCh: make(chan []byte),
	}, nil
}

func (d *DiskStorage) StartStorageRoutine(ctx context.Context) (chan []byte, error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-d.storageCh:
				if err := d.append(data); err != nil {
					d.log.Error("failed to append data to file: %v", err)
				}
			}
		}
	}()

	return d.storageCh, nil
}

func (d *DiskStorage) append(data []byte) error {
	info, err := d.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}

	if info.Size()+int64(len(data)) > int64(d.batchSize) {
		err = d.close()
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

func (d *DiskStorage) close() error {
	return d.file.Close()
}
