package actions

import (
	"github.com/rugwirobaker/paypack-backend/pkg/ussd"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/action"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/command"
)

type action1 struct {
	action.Base
}

// Action1 contains the logic to view user's associated propery codes
func Action1() action.Action {
	action := &action1{
		Base: action.NewBase(1),
	}
	// action.Register(1, initAction11())
	// action.Register(2, initAction12())
	return action
}

func (action *action1) Depth() int {
	return action.Base.Depth()
}

func (action *action1) Run(cmd *command.Command) (ussd.Result, error) {
	if cmd.Len() > action.Depth()+1 {
		next, err := action.Next(cmd)
		if err != nil {
			return nil, err
		}
		return next.Run(cmd)
	}
	return action.run()
}

func (action *action1) Next(cmd *command.Command) (action.Action, error) {
	return action.Base.Next(cmd)
}

func (action *action1) Register(key int, sub action.Action) {
	action.Base.Register(key, sub)
}

func (action *action1) Tail() bool {
	return action.Base.Tail()
}

func (action *action1) run() (ussd.Result, error) {
	var menu = "Kwemeza kureba code y'(z')inzu, Shyiramo nimero za telephone"
	return NewResult(menu), nil
}
