package batch

import (
	"testing"

	"github.com/matryer/is"
)

func Test(t *testing.T) {
	is := is.New(t)

	type r struct {
		start, end int
	}
	var ranges []r
	err := All(100, 10, func(start, end int) error {
		ranges = append(ranges, r{
			start: start,
			end:   end,
		})
		return nil
	})
	is.NoErr(err)

	is.Equal(len(ranges), 10)
	is.Equal(ranges[0].start, 0)
	is.Equal(ranges[0].end, 9)
	is.Equal(ranges[1].start, 10)
	is.Equal(ranges[1].end, 19)
	is.Equal(ranges[2].start, 20)
	is.Equal(ranges[2].end, 29)
	is.Equal(ranges[3].start, 30)
	is.Equal(ranges[3].end, 39)
	is.Equal(ranges[4].start, 40)
	is.Equal(ranges[4].end, 49)
	is.Equal(ranges[5].start, 50)
	is.Equal(ranges[5].end, 59)
	is.Equal(ranges[6].start, 60)
	is.Equal(ranges[6].end, 69)
	is.Equal(ranges[7].start, 70)
	is.Equal(ranges[7].end, 79)
	is.Equal(ranges[8].start, 80)
	is.Equal(ranges[8].end, 89)
	is.Equal(ranges[9].start, 90)
	is.Equal(ranges[9].end, 99)
}
