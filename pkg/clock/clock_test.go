package clock_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/clock"
	"github.com/stretchr/testify/assert"
)

func TestTimeIn(t *testing.T) {

	expected := time.Now()

	got := clock.TimeIn(expected.UTC(), clock.EAST)

	assert.True(t, expected.Local().Equal(got), fmt.Sprintf("expected %v got %v", expected, got))

}

func TestFormat(t *testing.T) {
	cases := []struct {
		input    time.Time
		expected string
		layout   string
	}{
		{
			input:    time.Date(2020, time.January, 1, 12, 0, 0, 0, time.Local),
			expected: "1 Jan 2020 12:00",
			layout:   clock.LayoutCustom,
		},
		{
			input:    time.Date(1996, time.August, 7, 2, 30, 0, 0, time.Local),
			expected: "7 Aug 1996 02:30",
			layout:   clock.LayoutCustom,
		},
	}

	for _, tc := range cases {
		got := clock.Format(tc.input, tc.layout)
		assert.Equal(t, tc.expected, got, fmt.Sprintf("expected '%s' got '%s'", tc.expected, got))
	}
}
