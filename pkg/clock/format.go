package clock

import "time"

// Time format layouts
const (
	LayoutCustom = "2 Jan 2006 15:04"
)

// Format time string
func Format(t time.Time, layout string) string {
	return t.Format(layout)
}
