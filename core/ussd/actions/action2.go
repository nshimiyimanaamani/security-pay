package actions

import (
	"github.com/rugwirobaker/paypack-backend/pkg/ussd"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/action"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/command"
)

type action2 struct {
	action.Base
	output string
}

// Action2 ...
func Action2() action.Action {
	action := &action2{
		output: "action:1",
		Base:   action.New(1),
	}
	// action.Register(1, initAction11())
	// action.Register(2, initAction12())
	return action
}

func (action *action2) Depth() int {
	return action.Base.Depth()
}

func (action *action2) Run(cmd *command.Command) (ussd.Result, error) {
	if cmd.Len() > action.Depth()+1 {
		next, err := action.Next(cmd)
		if err != nil {
			return nil, err
		}
		return next.Run(cmd)
	}
	return action.run()
}

func (action *action2) Next(cmd *command.Command) (action.Action, error) {
	return action.Base.Next(cmd)
}

func (action *action2) Register(key int, sub action.Action) {
	action.Base.Register(key, sub)
}

func (action *action2) Tail() bool {
	return action.Base.Tail()
}

func (action *action2) run() (ussd.Result, error) {
	var menu = "Kwishyura, Shyiramo code y' inzu"

	return NewResult(menu), nil
}
