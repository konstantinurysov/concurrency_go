package server

import (
	"fmt"
)

func (s *Server) readAndParseCommand() (string, []string, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	return s.parser.Parse(line)
}

func (s *Server) handleSet(args []string) string {
	if len(args) < 2 {
		return "failed to set: not enough arguments"
	}
	s.engine.Set(args[0], args[1])
	return "ok"
}

func (s *Server) handleGet(args []string) string {
	if len(args) < 1 {
		return "failed to get: not enough arguments"
	}
	if val, ok := s.engine.Get(args[0]); ok {
		return val
	}
	return " "
}

func (s *Server) handleDel(args []string) string {
	if len(args) < 1 {
		return "failed to delete: not enough arguments"
	}
	s.engine.Delete(args[0])

	return "ok"
}

func (s *Server) handleHelp(args []string) string {
	return guide
}

func (s *Server) dispatchCommand(command string, args []string) string {
	// if command == exitCommand {
	// 	return "", io.EOF
	// }
	if cmdDef, ok := s.commands[command]; ok {
		s.logger.Info("command: %s, args: %v", command, args)
		if len(args) < cmdDef.minArgs {
			return fmt.Sprintf("command %s requires at least %d argument(s), got %d", command, cmdDef.minArgs, len(args))
		}
		return cmdDef.handler(args)
	}
	return fmt.Sprintf("unknown command: %s", command)
}
