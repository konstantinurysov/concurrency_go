package network

import (
	"concurrency_hw1/internal/config"
	"concurrency_hw1/pkg/logger"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type TCPHandler = func(context.Context, []byte) []byte

type ServerInterface interface {
	Execute(ctx context.Context, handleRequest func(ctx context.Context, request []byte) []byte) error
}

type Server struct {
	tcpServer *TCPServer
	logger    logger.LoggerInterface
}

func NewServer(cfg *config.Config, logger logger.LoggerInterface) (*Server, error) {
	tcpServer, err := NewTCPServer(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create TCP server: %w", err)
	}
	return &Server{
		tcpServer: tcpServer,
		logger:    logger,
	}, nil
}

func (s *Server) Execute(ctx context.Context, handleRequest func(ctx context.Context, request []byte) []byte) error {
	wg := sync.WaitGroup{}

	fmt.Println("Welcome to SuperKV database. Waiting for your commands")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			c, err := s.tcpServer.listener.Accept()

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				s.logger.Error("failed to accept connection: %v", err.Error())
				return
			}

			go func() {
				s.handleConnection(ctx, c, handleRequest)
			}()
		}
	}()

	<-ctx.Done()
	s.logger.Info("stopping TCP server")
	if err := s.tcpServer.listener.Close(); err != nil {
		s.logger.Error("failed to close listener: %v", err.Error())
	}

	wg.Wait()

	return nil
}

func (s *Server) handleConnection(ctx context.Context, connection net.Conn, handler TCPHandler) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("captured panic", v)
		}

		if err := connection.Close(); err != nil {
			s.logger.Warn("failed to close connection: %v", err)
		}
	}()
	request := make([]byte, s.tcpServer.bufferSize)

Loop:
	for {
		if s.tcpServer.idleTimeout != 0 {
			if err := connection.SetReadDeadline(time.Now().Add(s.tcpServer.idleTimeout)); err != nil {
				s.logger.Warn("failed to set read deadline %v", err.Error())
				break Loop
			}
		}
		select {
		case <-ctx.Done():
			break Loop
		default:
			count, err := connection.Read(request)
			if err != nil && errors.Is(err, io.EOF) {
				s.logger.Warn("failed to read from connection: %v", err.Error())
				break Loop
			} else if count == s.tcpServer.bufferSize {
				s.logger.Warn("buffer size is too small")
				break Loop
			}

			if s.tcpServer.idleTimeout != 0 {
				if err := connection.SetWriteDeadline(time.Now().Add(s.tcpServer.idleTimeout)); err != nil {
					s.logger.Warn("failed to set read deadline %v", err.Error())
					break Loop
				}
			}
			s.logger.Info("request: %v", string(request))
			response := handler(ctx, request[:count])
			s.logger.Info("response: %v", string(response))
			if _, err := connection.Write(response); err != nil {
				s.logger.Warn(
					"failed to write data to %v: %v",
					connection.RemoteAddr().String(),
					err.Error(),
				)
				break Loop
			}
		}
	}

}
