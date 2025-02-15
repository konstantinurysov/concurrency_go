package server

import (
	"bufio"
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/config"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/logger"
	"errors"
	"net"
	"sync"
	"time"

	"context"
	"fmt"
	"io"
	"os"
)

const (
	setCommand  = "SET"
	getCommand  = "GET"
	delCommand  = "DEL"
	helpCommand = "help"
	exitCommand = "exit"
	guide       = "query = set_command | get_command | del_command \n set_command = \"SET\" argument argument \n get_command = \"GET\" argument \n del_command = \"DEL\" argument \n argument    = punctuation | letter | digit { punctuation | letter | digit } \n punctuation = \"*\" | \"/\" | \"_\" | ... \n letter      = \"a\" | ... | \"z\" | \"A\" | ... | \"Z\" \n digit       = \"0\" | ... | \"9\" \n exit_command = \"exit\""
)

type commandFunc func(args []string) string
type TCPHandler = func(context.Context, []byte) []byte

type CommandDefinition struct {
	minArgs int
	handler commandFunc
}

type Server struct {
	config    *config.Config
	logger    *logger.Logger
	reader    *bufio.Reader
	parser    *compute.Parser
	engine    *storage.Engine
	tcpServer *TCPServer
	commands  map[string]CommandDefinition
	mu        sync.Mutex
}

func NewServer(logger *logger.Logger, parser *compute.Parser, engine *storage.Engine, config *config.Config) *Server {
	tcpServer, err := NewTCPServer(config, logger)
	if err != nil {
		logger.Fatal("failed to create TCP server")
	}

	s := &Server{
		config:    config,
		logger:    logger,
		reader:    bufio.NewReader(os.Stdin),
		parser:    parser,
		engine:    engine,
		tcpServer: tcpServer,
		mu:        sync.Mutex{},
	}
	s.initCommands()

	return s
}

func (s *Server) initCommands() {
	s.commands = map[string]CommandDefinition{
		setCommand:  {minArgs: 2, handler: s.handleSet},
		getCommand:  {minArgs: 1, handler: s.handleGet},
		delCommand:  {minArgs: 1, handler: s.handleDel},
		helpCommand: {minArgs: 0, handler: s.handleHelp},
	}
}

func (s *Server) Execute(ctx context.Context) error {
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
			if s.increaseMaxConnections() {
				s.logger.Warn("max connections reached")
				c.Write([]byte("max connections reached please try again later"))
				continue
			}

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return nil
				}

				return fmt.Errorf("failed to accept connection: %w", err)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				s.handleConnection(ctx, c, s.handleRequest)
			}()
		}
	}
}

func (s *Server) increaseMaxConnections() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tcpServer.maxConnections == 0 {
		return true
	}
	s.tcpServer.maxConnections--
	return false
}

func (s *Server) handleRequest(ctx context.Context, request []byte) []byte {
	s.logger.Info("handleRequest request: %v", string(request))
	command, args, err := s.parser.Parse(string(request))
	if err != nil {
		s.logger.Error(err)
		return nil
	}

	return []byte(s.dispatchCommand(command, args))
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
