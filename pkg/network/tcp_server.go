package network

import (
	"concurrency_hw1/internal/config"
	"concurrency_hw1/pkg/common"
	"concurrency_hw1/pkg/logger"
	"fmt"

	"errors"
	"net"
	"time"
)

const defaultServerAddress = ":3223"

type TCPServer struct {
	listener net.Listener

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int
	logger         logger.LoggerInterface
}

func NewTCPServer(cfg *config.Config, logger logger.LoggerInterface) (*TCPServer, error) {
	if cfg == nil {
		logger.Fatal("empty config")
		return nil, nil
	}

	if logger == nil {
		logger.Fatal("empty logger")
	}

	listener, err := net.Listen("tcp", cfg.Network.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &TCPServer{
		listener: listener,
		logger:   logger,
	}

	options, err := server.getOptions(cfg.Network)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	for _, option := range options {
		option(server)
	}

	if server.bufferSize == 0 {
		server.bufferSize = 4 << 10
	}

	return server, nil

}

func (s *TCPServer) getOptions(cfg *config.NetworkConfig) ([]TCPServerOption, error) {
	var options []TCPServerOption

	if cfg.Address != "" {
		cfg.Address = defaultServerAddress
	}

	if cfg.MaxConnections != 0 {
		options = append(options, WithServerMaxConnectionsNumber(uint(cfg.MaxConnections)))
	}

	if cfg.MaxMessageSize != "" {
		size, err := common.ParseSize(cfg.MaxMessageSize)
		if err != nil {
			return nil, errors.New("incorrect max message size")
		}

		options = append(options, WithServerBufferSize(uint(size)))
	}

	if cfg.IdleTimeout != 0 {
		options = append(options, WithServerIdleTimeout(cfg.IdleTimeout))
	}

	return options, nil
}
