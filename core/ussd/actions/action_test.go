package actions

import (
	"testing"

	"github.com/rugwirobaker/paypack-backend/pkg/ussd/action"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/command"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/executor"
)

func newExecutor(root action.Action) executor.Executor {
	return executor.NewSimpleExecutor(root)
}

func TestAction(t *testing.T) {
	root := Action0()

	executor := newExecutor(root)

	cases := []struct {
		desc     string
		input    string
		expected string
		err      error
	}{
		{
			desc:     "action0: main menu",
			input:    "*662*102#",
			expected: "Murakaza neza kuri paypack\n1. reba code y' inzu yawe\n2. kwishyura\n",
			err:      nil,
		},
		{
			desc:     "action1: view properties",
			input:    "*662*102*1#",
			expected: "Kwemeza kureba code y'(z')inzu, Andika nimero za telephone",
			err:      nil,
		},
		{
			desc:     "action2: make payment",
			input:    "*662*102*2#",
			expected: "Kwishyura, Andika code y' inzu",
			err:      nil,
		},
	}

	for _, tc := range cases {
		cmd, err := command.Parse(tc.input)

		if err != tc.err {
			t.Fatalf("expected error: %v, got error %v", tc.err, err)
		}
		res, err := executor.Execute(cmd)

		if err != tc.err {
			t.Fatalf("expected error: %v, got error %v", tc.err, err)
		}

		got := res.String()
		if got != tc.expected {
			t.Errorf("expected: \"%s\" got: \"%s\"", tc.expected, got)
		}

	}
}
