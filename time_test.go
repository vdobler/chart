package chart

import (
	"testing"
	"time"
)

type roundingtest struct {
	date, expected string
	delta          TimeDelta
}

func TestRoundDown(t *testing.T) {
	samples := []roundingtest{{"2011-07-04 16:33:23", "2011-01-01 00:00:00", Year{1}},
		{"2011-07-04 16:33:23", "2010-01-01 00:00:00", Year{10}},
	}

	for _, sample := range samples {
		date, e1 := time.Parse("2006-01-02 15:04:05", sample.date)
		expected, e2 := time.Parse("2006-01-02 15:04:05", sample.expected)
		if e1 != nil || e2 != nil {
			t.FailNow()
		}
		sample.delta.RoundDown(date)
		if date.Seconds() != expected.Seconds() {
			t.Errorf("RoundDown %s to %s != %s, was %s", sample.date, sample.delta,
				sample.expected, date.Format("2006-01-02 15:04:05"))
		}
	}

}
