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


func TextCalendarWeek(t *testing.T) {
	for _, u := range [][4]int{
		[4]int{2011,1,1,52},
		[4]int{2011,1,2,52},
		[4]int{2011,1,3,1},
		[4]int{2011,1,4,1},
		[4]int{2011,1,5,1},
		[4]int{2011,1,6,1},
		[4]int{2011,1,7,1},
		[4]int{2011,1,8,1},
		[4]int{2011,1,9,1},
		[4]int{2011,1,10,2},
		[4]int{2011,12,25,51},
		[4]int{2011,12,26,52},
		[4]int{2011,12,27,52},
		[4]int{2011,12,28,52},
		[4]int{2011,12,29,52},
		[4]int{2011,12,30,52},
		[4]int{2011,12,31,52},
		[4]int{1995,1,1,52},
		[4]int{1995,1,2,1},
		[4]int{1996,1,1,1},
		[4]int{1996,1,7,1},
		[4]int{1996,1,8,2},
		[4]int{1997,1,1,1},
		[4]int{1998,1,1,1},
		[4]int{1999,1,1,53},
		[4]int{2000,1,1,52},
		[4]int{2001,1,1,1},
		[4]int{2002,1,1,1},
		[4]int{2003,1,1,1},
		[4]int{2004,1,1,1},
		[4]int{2005,1,1,53},
		[4]int{2006,1,1,52},
		[4]int{2007,1,1,1},
		[4]int{2008,1,1,1},
		[4]int{2009,1,1,1},
		[4]int{2010,1,1,53},
	} {
		dt := &time.Time{Year: int64(u[0]), Month: u[1], Day: u[2]}
		dt = time.SecondsToLocalTime(dt.Seconds())
		w := calendarWeek(dt)
		if w != u[3] {
			t.Errorf("Failed on %d-%d-%d. Got %d, expected %d.", u[0], u[1], u[2], w, u[3])
		}
	}

}