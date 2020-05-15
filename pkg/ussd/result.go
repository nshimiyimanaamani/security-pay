package ussd

import "fmt"

// Tailer ...
type Tailer interface {
	End() int
}

// Result ...
type Result interface {
	fmt.Stringer
	Tailer
}
