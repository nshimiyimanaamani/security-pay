package command

import "fmt"

// Type indicates instruction type which is either data or key
type Type int

// Types
const (
	Variable Type = iota
	Instruction
)

func (t Type) String() string {
	return []string{"Variable", "Instruction"}[t]
}

// Token corresponds to a single parsed instruction
type Token struct {
	Type  Type
	Value string
}

func (tk Token) String() string {
	return fmt.Sprintf("%s<%s>", tk.String(), tk.Value)
}
