package server

import (
	"fmt"
	"io"
)

func (s *Server) readAndParseCommand() (string, []string, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	return s.parser.Parse(line)
}

func (s *Server) handleSet(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("failed to set: not enough arguments")
	}
	s.engine.Set(args[0], args[1])
	return nil
}

func (s *Server) handleGet(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("failed to get: not enough arguments")
	}
	fmt.Println(s.engine.Get(args[0]))
	return nil
}

func (s *Server) handleDel(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("failed to delete: not enough arguments")
	}
	s.engine.Delete(args[0])

	return nil
}

func (s *Server) handleHelp(args []string) error {
	fmt.Println(guide)
	return nil
}

func (s *Server) dispatchCommand(command string, args []string) error {
	if command == exitCommand {
		return io.EOF
	}
	if cmdDef, ok := s.commands[command]; ok {
		if len(args) < cmdDef.minArgs {
			return fmt.Errorf("command %s requires at least %d argument(s), got %d", command, cmdDef.minArgs, len(args))
		}
		return cmdDef.handler(args)
	}
	return fmt.Errorf("unknown command: %s", command)
}
