package main

import (
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/config"
	"concurrency_hw1/internal/server"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/internal/wal"
	"concurrency_hw1/pkg/disk"
	"concurrency_hw1/pkg/logger"
	"flag"
	"os"
	"time"

	"context"
	"os/signal"
	"syscall"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	logger := logger.New("debug", "local")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	address := flag.String("address", "127.0.0.1:3223", "Address of the server")
	idleTimeoutStr := flag.String("idle_timeout", "5m", "Idle timeout for connection")
	walPath := flag.String("wal_path", "./wal", "Path to write ahead log")
	idleTimeout, err := time.ParseDuration(*idleTimeoutStr)
	if err != nil {
		logger.Fatal("failed to parse idle timeout", err)
	}
	maxConnections := flag.Int("max_connections", 100, "Max connections for server")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	if ConfigFileName == "" {
		ConfigFileName = "./../../config.yml"
	}

	cfg, err := config.Load(logger, ConfigFileName)
	if err != nil {
		logger.Info("failed to load config. working with default values")
		cfg = &config.Config{Network: &config.NetworkConfig{}}
		if address != nil {
			cfg.Network.Address = *address
		}
		if maxConnections != nil {
			cfg.Network.MaxConnections = *maxConnections
		}
		if maxMessageSizeStr != nil {
			cfg.Network.MaxMessageSize = *maxMessageSizeStr
		}
		cfg.Network.IdleTimeout = idleTimeout
		if walPath != nil {
			cfg.Storage.Path = *walPath
		}
	}

	parser := compute.NewParser()
	engine := storage.NewEngine()
	diskStorage, err := disk.NewDiskStorage(cfg.Storage.Path, cfg.Storage.MaxSegmentSize, logger)
	if err != nil {
		logger.Error("failed to create disk storage: %w", err)
		return
	}

	wal := wal.NewWALService(diskStorage, cfg.Storage.FlushingBatchSize, cfg.Storage.FlushingBatchTimeout, logger)
	wal.Start(ctx)
	service := server.NewServer(logger, parser, engine, wal.WALChannel, cfg)
	service.Execute(ctx)

	logger.Info("all services are stopped")
}
