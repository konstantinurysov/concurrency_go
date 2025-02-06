package main

import (
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/server"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/logger"

	"context"
	"os/signal"
	"syscall"
)

func main() {
	logger := logger.New("debug", "local")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	parser := compute.NewParser()
	engine := storage.NewEngine()

	service := server.NewServer(logger, parser, engine)
	service.Execute(ctx)

	logger.Info("all services are stopped")
}
