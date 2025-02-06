package server

import (
	"bufio"
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/logger"
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

type commandFunc func(args []string) error

type CommandDefinition struct {
	minArgs int
	handler commandFunc
}

type Server struct {
	logger   *logger.Logger
	reader   *bufio.Reader
	parser   *compute.Parser
	engine   *storage.Engine
	commands map[string]CommandDefinition
}

func NewServer(logger *logger.Logger, parser *compute.Parser, engine *storage.Engine) *Server {
	s := &Server{
		logger: logger,
		reader: bufio.NewReader(os.Stdin),
		parser: parser,
		engine: engine,
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
	fmt.Println("Welcome to SuperKV database. Waiting for your commands")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			command, args, err := s.readAndParseCommand()
			if err != nil {
				return err
			}

			if err := s.dispatchCommand(command, args); err != nil {
				if err == io.EOF {
					return nil
				}
				s.logger.Error(err)
			}
		}
	}
}
