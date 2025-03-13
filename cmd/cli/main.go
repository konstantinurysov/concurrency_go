package main

import (
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/config"
	"concurrency_hw1/internal/server"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/internal/wal"
	"concurrency_hw1/pkg/disk"
	"concurrency_hw1/pkg/logger"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	logger := logger.New("debug", "local")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	address := flag.String("address", "localhost:3223", "Address of the server")
	idleTimeout := flag.Duration("idle_timeout", 1, "Idle timeout for connection")
	maxConnections := flag.Int("max_connections", 100, "Max connections for server")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	if ConfigFileName == "" {
		ConfigFileName = "./../../config.yml"
	}

	cfg, err := config.Load(logger, ConfigFileName)
	if err != nil {
		if cfg.Network.Address == "" {
			cfg.Network.Address = *address
		}
		if cfg.Network.MaxConnections == 0 {
			cfg.Network.MaxConnections = *maxConnections
		}
		if cfg.Network.MaxMessageSize == "" {
			cfg.Network.MaxMessageSize = *maxMessageSizeStr
		}
		if cfg.Network.IdleTimeout == 0 {
			cfg.Network.IdleTimeout = *idleTimeout
		}
	}

	diskStorage, err := disk.NewDiskStorage(cfg.Storage.Path, cfg.Storage.MaxSegmentSize, logger)
	if err != nil {
		logger.Error("failed to create disk storage: %w", err)
		return
	}

	parser := compute.NewParser()
	engine := storage.NewEngine()
	wal := wal.NewWALService(diskStorage, cfg.Storage.FlushingBatchSize, cfg.Storage.FlushingBatchTimeout, logger)
	service := server.NewServer(logger, parser, engine, wal.WALChannel, cfg)
	service.Execute(ctx)

	logger.Info("all services are stopped")
}
