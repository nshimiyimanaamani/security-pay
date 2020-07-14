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

// End converts true to 0 and false to 1
func (res Result) End() int {
	if res.end {
		return 0
	}
	return 1
}
