package compute

import (
	"concurrency_hw1/internal/storage"
	"strings"
)

type CommandHandler interface {
	CanHandle(command string, args []string) bool
	Execute(args []string, engine *storage.Engine) (string, error)
}

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(line string) (string, []string, error) {
	if !p.Validate(line) {
		return "", nil, nil
	}

	arr := strings.Fields(line)
	if len(line) != 0 {
		return arr[0], arr[1:], nil
	}

	return "", nil, nil
}

// this function is used to validate input string
func (p *Parser) Validate(line string) bool {
	arr := strings.Fields(line)
	return len(arr) != 0
}
