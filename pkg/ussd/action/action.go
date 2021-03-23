package action

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rugwirobaker/paypack-backend/pkg/ussd"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/command"
)

var _ Action = (*Base)(nil)

// Action defines a ussd action
type Action interface {
	//Depth returns the depth at which the action is registered
	Depth() int

	// Run the command and return a result
	Run(*command.Command) (ussd.Result, error)

	// Next selects the next instruction in a command
	Next(*command.Command) (Action, error)

	// Register the action to a parent
	Register(int, Action)

	// Tail determines whether the action is the last in a session
	Tail() bool
}

// Base action
type Base struct {
	depth int
	subs  map[int]Action
}

// New ...
func New(depth int) Base {
	return Base{
		depth: depth,
		subs:  make(map[int]Action),
	}
}

// Depth ...
func (action *Base) Depth() int {
	return action.depth
}

// Run ...
func (action *Base) Run(*command.Command) (ussd.Result, error) {
	return nil, errors.New("not implemented")
}

// Next ...
func (action *Base) Next(cmd *command.Command) (Action, error) {
	token := cmd.Next(action.Depth())

	switch token.Type {
	case command.Instruction:
		key, err := strconv.Atoi(token.Value)
		if err != nil {
			return nil, fmt.Errorf("'%s'", err)
		}
		return action.subs[key], nil
	case command.Variable:
		return nil, fmt.Errorf("token '%s' is not a instruction", token)
	default:
		return nil, errors.New("action not found")
	}
}

// Register ...
func (action *Base) Register(key int, sub Action) {
	action.subs[key] = sub
}

// Tail indicates whether action is a tail
func (action *Base) Tail() bool {
	if len(action.subs) == 0 {
		return true
	}
	return false
}

//important: how does a command look like
