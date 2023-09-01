package actions

import (
	"github.com/nshimiyimanaamani/paypack-backend/pkg/ussd"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/ussd/action"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/ussd/command"
)

var _ action.Action = (*action0)(nil)

type action0 struct {
	action.Base
}

func (action *action0) Depth() int {
	return action.Base.Depth()
}

func (action *action0) Run(cmd *command.Command) (ussd.Result, error) {

	if cmd.Len() > action.Depth() {
		next, err := action.Next(cmd)
		if err != nil {
			return nil, err
		}
		return next.Run(cmd)
	}

	return action.run()
}

func (action *action0) Next(cmd *command.Command) (action.Action, error) {
	return action.Base.Next(cmd)
}

func (action *action0) Register(key int, sub action.Action) {
	action.Base.Register(key, sub)
}

func (action *action0) Tail() bool {
	return action.Base.Tail()
}

func (action *action0) run() (ussd.Result, error) {
	var menu = "Murakaza neza kuri paypack\n1. reba code y' inzu yawe\n2. kwishyura\n"

	return NewResult(menu), nil
}

// Action0 ...
func Action0() action.Action {
	var root = &action0{
		Base: action.New(0),
	}
	root.Register(1, Action1())
	root.Register(2, Action2())

	return root
}
