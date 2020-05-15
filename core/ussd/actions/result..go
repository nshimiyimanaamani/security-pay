package actions

// Result ...
type Result struct {
	output string
	end    bool
}

// NewResult ...
func NewResult(out string) *Result {
	return &Result{output: out}
}

func (res Result) String() string {
	return res.output
}

// End converts true to 1 and false to 0
func (res Result) End() int {
	if res.end {
		return 1
	}
	return 0
}
