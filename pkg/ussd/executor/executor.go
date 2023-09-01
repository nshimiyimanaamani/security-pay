package executor

import (
	"github.com/nshimiyimanaamani/paypack-backend/pkg/ussd"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/ussd/action"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/ussd/command"
)

// Executor ...
type Executor interface {
	Execute(*command.Command) (ussd.Result, error)
}

var _ Executor = (*simpleExecutor)(nil)

type simpleExecutor struct {
	base action.Action
}

// NewSimpleExecutor ...
func NewSimpleExecutor(base action.Action) Executor {
	return &simpleExecutor{
		base: base,
	}
}

func (exec *simpleExecutor) Execute(cmd *command.Command) (ussd.Result, error) {
	cmd.Skip(1)
	return exec.base.Run(cmd)
}
