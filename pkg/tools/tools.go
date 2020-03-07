package tools

import "time"

// LastMonth ...
func LastMonth() time.Time {
	now := time.Now()
	return now.AddDate(0, -1, 0)
}

// Yesterday ...
func Yesterday() time.Time {
	now := time.Now()
	return now.AddDate(0, 0, -1)
}

// AddMonth ...
func AddMonth(t time.Time, i int) time.Time {
	return t.AddDate(0, i, 0)
}
