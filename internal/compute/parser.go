package compute

import (
	"fmt"
	"strings"
)

type Parser struct {
}

func (p *Parser) Parse(line string) (string, []string, error) {
	if line == "\n" {
		return "", nil, nil
	}

	if len(line) < 3 {
		err := fmt.Errorf("failed to parse command: too short")
		return "", nil, err
	}

	arr := strings.Fields(line)
	if len(arr) < 2 && arr[0] != "help" {
		err := fmt.Errorf("failed to parse command: not enough arguments")
		return "", nil, err
	}

	return arr[0], arr[1:], nil
}

func NewParser() *Parser {
	return &Parser{}
}
