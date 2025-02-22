package server

import (
	"bufio"
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/config"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/concurrency"
	"concurrency_hw1/pkg/logger"
	"concurrency_hw1/pkg/network"

	"context"
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
	server    *network.Server
	commands  map[string]CommandDefinition
	semaphore *concurrency.Semaphore
}

func NewServer(logger *logger.Logger, parser *compute.Parser, engine *storage.Engine, config *config.Config) *Server {
	server, err := network.NewServer(config, logger)
	if err != nil {
		logger.Fatal(err)
	}

	s := &Server{
		config:    config,
		logger:    logger,
		reader:    bufio.NewReader(os.Stdin),
		parser:    parser,
		engine:    engine,
		server:    server,
		semaphore: concurrency.NewSemaphore(config.Network.MaxConnections),
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
	s.server.Execute(ctx, s.handleRequest)
	return nil
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
