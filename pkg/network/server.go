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

type Server struct {
	tcpServer *TCPServer
	logger    *logger.Logger
}

func NewServer(cfg *config.Config, logger *logger.Logger) (*Server, error) {
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
	defer func() {
		wg.Wait()
		s.tcpServer.listener.Close()
	}()

	fmt.Println("Welcome to SuperKV database. Waiting for your commands")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			c, err := s.tcpServer.listener.Accept()

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return nil
				}

				return fmt.Errorf("failed to accept connection: %w", err)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				s.handleConnection(ctx, c, handleRequest)
			}()
		}
	}
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
