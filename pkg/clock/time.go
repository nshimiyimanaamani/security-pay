package clock

import "time"

// Zone typfies a single timezone
type Zone string

// timezones
const (
	Rwanda Zone = "rwanda"
)

// Zones defines time zones
var zones = map[Zone]string{
	Rwanda: "Africa/Kigali",
}

// TimeIn returns given time set to given location
func TimeIn(t time.Time, zone Zone) time.Time {
	name := zones[zone]

	loc, _ := time.LoadLocation(name)

	return t.In(loc)
}
