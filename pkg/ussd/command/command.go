package command

import (
	"strings"
)

// Command contains a list of instructions
type Command struct {
	Tokens []Token
}

// Next returns the next instruction
func (cmd *Command) Next(prev int) *Token {
	if prev > cmd.Len() {
		return nil
	}
	return &cmd.Tokens[prev]
}

// Push appends token to Command
func (cmd *Command) Push(token Token) {
	cmd.Tokens = append(cmd.Tokens, token)
}

// Skip skips instructions starting from index
func (cmd *Command) Skip(index int) {
	cmd.Tokens = append([]Token{}, cmd.Tokens[index+1:]...)
}

// Len Returns the length of tokens
func (cmd *Command) Len() int {
	return len(cmd.Tokens)
}

// Parse parses user input into executable command
func Parse(in string) (*Command, error) {

	var tokens []Token

	in = strings.Trim(in, "#")

	values := strings.Split(in, "*")

	for _, value := range values {
		var token Token

		if value == "" {
			continue
		}

		if len(value) <= 3 {
			token = newInstruction(value)
			tokens = append(tokens, token)
		} else {
			token = newVariable(value)
			tokens = append(tokens, token)
		}
	}

	return &Command{Tokens: tokens}, nil
}

func newVariable(value string) Token {
	return Token{
		Type:  Variable,
		Value: value,
	}
}

func newInstruction(value string) Token {
	return Token{
		Type:  Instruction,
		Value: value,
	}
}
