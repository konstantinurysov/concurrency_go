package main

import (
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/server"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/logger"
	"concurrency_hw1/pkg/signal"
	"context"
	"fmt"

	"github.com/oklog/run"
)

func main() {
	logger := logger.New("debug", "local")
	ctx, cancel := context.WithCancel(context.Background())

	parser := compute.NewParser()
	engine := storage.NewEngine()

	group := &run.Group{}
	{
		service := signal.NewService(cancel)
		group.Add(service.Execute, service.Interrupt)
	}
	{
		service := server.NewServer(ctx, cancel, logger, parser, engine)
		group.Add(service.Execute, service.Interrupt)
	}

	if err := group.Run(); err != nil {
		logger.Fatal(fmt.Errorf("the service stopped because of an error: %w", err))
	}

	logger.Info("all services are stopped")
}
