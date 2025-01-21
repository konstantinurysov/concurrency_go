package server

import (
	"bufio"
	"concurrency_hw1/internal/compute"
	"concurrency_hw1/internal/storage"
	"concurrency_hw1/pkg/logger"
	"context"
	"fmt"
	"os"
)

const (
	setCommand  = "SET"
	getCommand  = "GET"
	delCommand  = "DEL"
	helpCommand = "help"
	exitCommand = "exit"
	guide       = "query = set_command | get_command | del_command \n set_command = \"SET\" argument argument \n get_command = \"GET\" argument \n del_command = \"DEL\" argument \n argument    = punctuation | letter | digit { punctuation | letter | digit } \n punctuation = \"*\" | \"/\" | \"_\" | ... \n letter      = \"a\" | ... | \"z\" | \"A\" | ... | \"Z\" \n digit       = \"0\" | ... | \"9\" "
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *logger.Logger
	reader *bufio.Reader
	parser *compute.Parser
	engine *storage.Engine
}

func (s *Server) Execute() error {
	fmt.Println("Welcome to SuperKV database ) Waiting for your commands")
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			line, err := s.reader.ReadString('\n')
			if err != nil {
				return err
			}
			command, args, err := s.parser.Parse(line)
			if err != nil {
				s.logger.Error(err)
			}

			switch command {
			case setCommand:
				if len(args) < 2 {
					s.logger.Error("failed to set: not enough arguments")
				} else {
					s.engine.Set(args[0], args[1])
				}
			case getCommand:
				if len(args) < 1 {
					s.logger.Error("failed to set: not enough arguments")
				} else {
					fmt.Printf("%s\n", s.engine.Get(args[0]))
				}
			case delCommand:
				if len(args) < 1 {
					s.logger.Error("failed to delete: not enough arguments")
				} else {
					s.engine.Delete(args[0])
					fmt.Printf("%s removed\n", args[0])
				}
			case helpCommand:
				fmt.Println(guide)
			case exitCommand:
				return nil
			default:
				s.logger.Error("wrong command, please use help command")
			}
		}
	}
}

func (s *Server) Interrupt(err error) {
	s.cancel()
}

func NewServer(ctx context.Context, cancel context.CancelFunc, logger *logger.Logger, parser *compute.Parser, engine *storage.Engine) *Server {
	return &Server{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		reader: bufio.NewReader(os.Stdin),
		parser: parser,
		engine: engine,
	}
}
