package clock

import "time"

// Zone typfies a single timezone
type Zone string

// timezones
const (
	EAST Zone = "east-africa-standard-time"
)

// Zones defines time zones
var zones = map[Zone]string{
	EAST: "Africa/Nairobi",
}

// TimeIn returns given time set to given location
func TimeIn(t time.Time, zone Zone) time.Time {
	name := zones[zone]

	loc, err := time.LoadLocation(name)
	if err != nil {
		return t
	}

	return t.In(loc)
}
