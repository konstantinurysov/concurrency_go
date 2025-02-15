package main

import (
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/config"
	"concurrency_hw1/internal/server"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/logger"
	"flag"
	"os"

	"context"
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
		ConfigFileName = "../../config.yml"
	}

	cfg, err := config.Load(ConfigFileName, *address, *idleTimeout, *maxConnections, *maxMessageSizeStr)
	if err != nil {
		logger.Fatal(err)
	}

	parser := compute.NewParser()
	engine := storage.NewEngine()

	service := server.NewServer(logger, parser, engine, cfg)
	service.Execute(ctx)

	logger.Info("all services are stopped")
}
