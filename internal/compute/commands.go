package compute

import "concurrency_hw1/internal/storage"

const (
	setCommand      = "SET"
	setCommandSize  = 3
	getCommand      = "GET"
	getCommandSize  = 2
	delCommand      = "DEL"
	delCommandSize  = 2
	helpCommand     = "help"
	helpCommandSize = 1
	exitCommand     = "exit"
	exitCommandSize = 1
	guide           = "query = set_command | get_command | del_command \n set_command = \"SET\" argument argument \n get_command = \"GET\" argument \n del_command = \"DEL\" argument \n argument    = punctuation | letter | digit { punctuation | letter | digit } \n punctuation = \"*\" | \"/\" | \"_\" | ... \n letter      = \"a\" | ... | \"z\" | \"A\" | ... | \"Z\" \n digit       = \"0\" | ... | \"9\" "
)

type SetCommand struct {
}

func (s SetCommand) CanHandle(command string, args []string) bool {
	return command == setCommand && len(args) == setCommandSize
}

func (s SetCommand) Execute(args []string, engine *storage.Engine) (string, error) {
	engine.Set(args[0], args[1])
	return "", nil
}
