package chart

import (
	"fmt"
	"time"
)

// Represents a tic-distance in a timed axis
type TimeDelta interface {
	Seconds() int64 // amount of delta in seconds
	RoundDown(t *time.Time)
	String() string
	Format(t *time.Time) string
	Period() bool
}


// Second
type Second struct {
	Num int
}

func (s Second) Seconds() int64 { return int64(s.Num) }
func (s Second) RoundDown(t *time.Time) {
	t.Second = s.Num * (t.Second / s.Num)
}
func (s Second) String() string             { return fmt.Sprintf("%d seconds(s)", s.Num) }
func (s Second) Format(t *time.Time) string { return fmt.Sprintf(":%02d:%02d", t.Minute, t.Second) }
func (s Second) Period() bool               { return false }

// Minute
type Minute struct {
	Num int
}

func (m Minute) Seconds() int64 { return int64(60 * m.Num) }
func (m Minute) RoundDown(t *time.Time) {
	t.Second = 0
	t.Minute = m.Num * (t.Minute / m.Num)
}
func (m Minute) String() string             { return fmt.Sprintf("%d minute(s)", m.Num) }
func (m Minute) Format(t *time.Time) string { return fmt.Sprintf(":%02d", t.Minute) }
func (m Minute) Period() bool               { return false }

// Hour
type Hour struct {
	Num int
}

func (h Hour) Seconds() int64 { return 60 * 60 * int64(h.Num) }
func (h Hour) RoundDown(t *time.Time) {
	t.Second = 0
	t.Minute = 0
	t.Hour = h.Num * (t.Hour / h.Num)
}
func (h Hour) String() string             { return fmt.Sprintf("%d hours(s)", h.Num) }
func (h Hour) Format(t *time.Time) string { return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute) }
func (h Hour) Period() bool               { return false }

// Day
type Day struct {
	Num int
}

func (d Day) Seconds() int64 { return 60 * 60 * 24 * int64(d.Num) }
func (d Day) RoundDown(t *time.Time) {
	t.Hour, t.Minute, t.Second = 0, 0, 0
	t.Day = d.Num*((t.Day-1)/d.Num) + 1
	// t = time.SecondsToLocalTime(t.Second() + 60*60*2 ) // 2 hours 
	// y.RoundDown(t)
}
func (d Day) String() string             { return fmt.Sprintf("%d day(s)", d.Num) }
func (d Day) Format(t *time.Time) string { return fmt.Sprintf("%s", t.Format("Mon")) }
func (d Day) Period() bool               { return true }

// Week
type Week struct {
	Num int
}

func (w Week) Seconds() int64 { return 60 * 60 * 24 * 7 * int64(w.Num) }
func (w Week) RoundDown(t *time.Time) {
	shift := int64(60 * 60 * 24 * (time.Monday - t.Weekday))
	t.Hour, t.Minute, t.Second = 12, 0, 0 // Safeguard shift below againts daylightsavings and that like
	t = time.SecondsToLocalTime(t.Seconds() - shift)
	t.Hour, t.Minute, t.Second = 0, 0, 0
}
func (w Week) String() string { return fmt.Sprintf("%d week(s)", w.Num) }
func (w Week) Format(t *time.Time) string {
	// TODO(vodo): check if suitable
	jan01 := *t
	jan01.Month, jan01.Day, jan01.Hour, jan01.Minute, jan01.Second = 1, 1, 0, 0, 0
	diff := t.Seconds() - jan01.Seconds()
	week := int(float64(diff)/float64(60*60*24*7) + 0.5)
	if week == 0 {
		week++
	}
	return fmt.Sprintf("KW %d", week)
}
func (w Week) Period() bool { return true }

// Month
type Month struct {
	Num int
}

func (m Month) Seconds() int64 { return 60 * 60 * 24 * 365.25 / 12 * int64(m.Num) }
func (m Month) RoundDown(t *time.Time) {
	t.Day, t.Hour, t.Minute, t.Second = 1, 0, 0, 0
	t.Month = m.Num*((t.Month-1)/m.Num) + 1
	// weird looking "add some stuf, reparse and truncate again" fixes timezone shifts
	t = time.SecondsToLocalTime(t.Seconds() + 60*60*24) // 1 days
	t.Day, t.Hour, t.Minute, t.Second = 1, 0, 0, 0
}
func (m Month) String() string             { return fmt.Sprintf("%d month(s)", m.Num) }
func (m Month) Format(t *time.Time) string { return fmt.Sprintf("%02d.%d", t.Month, t.Year) }
func (m Month) Period() bool               { return true }

// Year
type Year struct {
	Num int
}

func (y Year) Seconds() int64 { return 60 * 60 * 24 * 365.25 * int64(y.Num) }
func (y Year) RoundDown(t *time.Time) {
	t.Hour, t.Minute, t.Second = 0, 0, 0
	t.Month, t.Day = 1, 1
	t.Year = int64(y.Num) * (t.Year / int64(y.Num))
	t = time.SecondsToLocalTime(t.Seconds() + 60*60*24*5) // 5 days
	t.Hour, t.Minute, t.Second = 0, 0, 0
	t.Month, t.Day = 1, 1
}
func (y Year) String() string             { return fmt.Sprintf("%d year(s)", y.Num) }
func (y Year) Format(t *time.Time) string { return fmt.Sprintf("%d", t.Year) }
func (y Year) Period() bool               { return true }


// The time deltas to use in 
// must be sorted min to max according to Seconds() of each member
var Delta []TimeDelta = []TimeDelta{Second{1}, Second{5}, Second{15},
	Minute{1}, Minute{5}, Minute{15},
	Hour{1}, Hour{6}, Day{1}, Week{1},
	Month{1}, Month{3}, Month{6},
	Year{1}, Year{10},
}


func RoundUp(tp *time.Time, d TimeDelta) *time.Time {
	// works only because all TimeDeltas are more than 3 times as large as the next lower
	tc := *tp
	d.RoundDown(&tc)
	shift := d.Seconds()
	shift += shift / 2
	t := time.SecondsToLocalTime(tc.Seconds() + shift)
	d.RoundDown(t)
	return t
}

func RoundNext(t *time.Time, d TimeDelta) *time.Time {
	os := t.Seconds()
	lt := *t
	d.RoundDown(&lt)
	shift := d.Seconds()
	ut := time.SecondsToLocalTime(lt.Seconds() + shift + shift/2) // see RoundUp()
	d.RoundDown(ut)
	ld := os - lt.Seconds()
	ud := ut.Seconds() - os
	if ld < ud {
		return &lt
	}
	return ut
}

func RoundDown(tp *time.Time, d TimeDelta) *time.Time {
	tc := *tp
	d.RoundDown(&tc)
	// fmt.Printf("RoundDown( %s )  -->  %s\n", tp.Format("2006-01-02 15:04:05 (Mon)"), 
	//	tc.Format("2006-01-02 15:04:05 (Mon)"))
	return &tc
}


func NextTimeDelta(d TimeDelta) TimeDelta {
	var i = 0
	sec := d.Seconds()
	for i < len(Delta) && Delta[i].Seconds() <= sec {
		i++
	}
	if i < len(Delta) {
		return Delta[i]
	}
	return Delta[len(Delta)-1]
}

func MatchingTimeDelta(delta float64, fac float64) TimeDelta {
	var i = 0
	for i+1 < len(Delta) && delta > fac*float64(Delta[i+1].Seconds()) {
		i++
	}
	// fmt.Printf("MatchingTimeDelta(%g): i=%d, %s...%s  ==  %d...%d\n  %t\n", delta, i, Delta[i], Delta[i+1], Delta[i].Seconds(), Delta[i+1].Seconds(), i+1 < len(Delta) && delta > fac * float64(Delta[i+1].Seconds()))
	if i+1 < len(Delta) {
		return Delta[i+1]
	}
	return Delta[len(Delta)-1]
}
