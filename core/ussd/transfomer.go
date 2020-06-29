package ussd

// Transformer turns cmd.Pattterns intoa acceptable inputs
type Transformer func(pattern string) string

func nopTransformer(pattern string) string {
	return pattern
}
