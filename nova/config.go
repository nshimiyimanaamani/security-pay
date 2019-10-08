package nova

import "time"

// Config contains the novapay API configuration
type Config struct {
	Endpoint string
	Token    string
	TimeOut  time.Duration
}
