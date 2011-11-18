package chart

import (
	"fmt"
	"time"
)

// Represents a tic-distance in a timed axis
type TimeDelta interface {
	Seconds() int64             // amount of delta in seconds
	RoundDown(t *time.Time)     // Round dow t to "whole" delta
	String() string             // retrieve string representation
	Format(t *time.Time) string // format t properly
	Period() bool               // true if this delta is a time period (like a month)
}

// Copy value of src to dest.
func cpTime(dest, src *time.Time) {
	dest.Year, dest.Month, dest.Day = src.Year, src.Month, src.Day
	dest.Hour, dest.Minute, dest.Second = src.Hour, src.Minute, src.Second
	dest.Weekday, dest.ZoneOffset, dest.Zone = src.Weekday, src.ZoneOffset, src.Zone
}

// Second
type Second struct {
	Num int
}

func (s Second) Seconds() int64             { return int64(s.Num) }
func (s Second) RoundDown(t *time.Time)     { t.Second = s.Num * (t.Second / s.Num) }
func (s Second) String() string             { return fmt.Sprintf("%d seconds(s)", s.Num) }
func (s Second) Format(t *time.Time) string { return fmt.Sprintf("%02d'%02d\"", t.Minute, t.Second) }
func (s Second) Period() bool               { return false }

// Minute
type Minute struct {
	Num int
}

func (m Minute) Seconds() int64             { return int64(60 * m.Num) }
func (m Minute) RoundDown(t *time.Time)     { t.Second, t.Minute = 0, m.Num*(t.Minute/m.Num) }
func (m Minute) String() string             { return fmt.Sprintf("%d minute(s)", m.Num) }
func (m Minute) Format(t *time.Time) string { return fmt.Sprintf("%02d'", t.Minute) }
func (m Minute) Period() bool               { return false }

// Hour
type Hour struct{ Num int }

func (h Hour) Seconds() int64             { return 60 * 60 * int64(h.Num) }
func (h Hour) RoundDown(t *time.Time)     { t.Second, t.Minute, t.Hour = 0, 0, h.Num*(t.Hour/h.Num) }
func (h Hour) String() string             { return fmt.Sprintf("%d hours(s)", h.Num) }
func (h Hour) Format(t *time.Time) string { return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute) }
func (h Hour) Period() bool               { return false }

// Day
type Day struct{ Num int }

func (d Day) Seconds() int64 { return 60 * 60 * 24 * int64(d.Num) }
func (d Day) RoundDown(t *time.Time) {
	t.Day, t.Hour, t.Minute, t.Second = d.Num*((t.Day-1)/d.Num)+1, 0, 0, 0
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
	org := t.Format("Mon 2006-01-02")
	week := calendarWeek(t)
	shift := int64(60 * 60 * 24 * (t.Weekday - time.Monday))
	t.Hour, t.Minute, t.Second = 0, 0, 0
	cpTime(t, time.SecondsToLocalTime(t.Seconds()-shift))
	t.Hour, t.Minute, t.Second = 0, 0, 0

	// daylight saving and that like might lead to different real shift
	for calendarWeek(t) < week {
		cpTime(t, time.SecondsToLocalTime(t.Seconds()+60*60*36))
		t.Hour, t.Minute, t.Second = 0, 0, 0
	}
	for calendarWeek(t) > week {
		cpTime(t, time.SecondsToLocalTime(t.Seconds()-60*60*36))
		t.Hour, t.Minute, t.Second = 0, 0, 0
	}
	trace.Printf("Week.Roundown(%s) --> %s", org, t.Format("Mon 2006-01-02"))

}
func (w Week) String() string             { return fmt.Sprintf("%d week(s)", w.Num) }
func (w Week) Format(t *time.Time) string { return fmt.Sprintf("W %d", calendarWeek(t)) }
func (w Week) Period() bool               { return true }

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
func (m Month) String() string { return fmt.Sprintf("%d month(s)", m.Num) }
func (m Month) Format(t *time.Time) string {
	if m.Num == 3 { // quarter years
		return fmt.Sprintf("Q%d %d", (t.Month-1)/3+1, t.Year)
	}
	if m.Num == 6 { // half years
		return fmt.Sprintf("H%d %d", (t.Month-1)/6+1, t.Year)
	}
	return fmt.Sprintf("%02d.%d", t.Month, t.Year)
}
func (m Month) Period() bool { return true }

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

// Delta is a list of increasing time deltas used to construct tic spacings
// for date/time axis.
// Must be sorted min to max according to Seconds() of each member.
var Delta []TimeDelta = []TimeDelta{
	Second{1}, Second{5}, Second{15},
	Minute{1}, Minute{5}, Minute{15},
	Hour{1}, Hour{6},
	Day{1}, Week{1},
	Month{1}, Month{3}, Month{6},
	Year{1}, Year{10},
}

// RoundUp will round tp up to next "full" d.
func RoundUp(tp *time.Time, d TimeDelta) *time.Time {
	// works only because all TimeDeltas are more than 1.5 times as large as the next lower
	tc := *tp
	d.RoundDown(&tc)
	shift := d.Seconds()
	shift += shift / 2
	t := time.SecondsToLocalTime(tc.Seconds() + shift)
	d.RoundDown(t)
	trace.Printf("RoundUp( %s, %s ) --> %s ", tp.Format("2006-01-02 15:04:05 (Mon)"), d.String(),
		t.Format("2006-01-02 15:04:05 (Mon)"))
	return t
}

// RoundNext will round t to nearest full d.
func RoundNext(t *time.Time, d TimeDelta) *time.Time {
	trace.Printf("RoundNext( %s, %s )", t.Format("2006-01-02 15:04:05 (Mon)"), d.String())
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

// RoundDown will round tp down to next "full" d.
func RoundDown(tp *time.Time, d TimeDelta) *time.Time {
	tc := *tp
	d.RoundDown(&tc)
	trace.Printf("RoundDown( %s, %s ) --> %s", tp.Format("2006-01-02 15:04:05 (Mon)"), d.String(),
		tc.Format("2006-01-02 15:04:05 (Mon)"))
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
	trace.Printf("MatchingTimeDelta(%g): i=%d, %s...%s  ==  %d...%d\n  %t\n",
		delta, i, Delta[i], Delta[i+1], Delta[i].Seconds(), Delta[i+1].Seconds(),
		i+1 < len(Delta) && delta > fac*float64(Delta[i+1].Seconds()))
	if i+1 < len(Delta) {
		return Delta[i+1]
	}
	return Delta[len(Delta)-1]
}

func dayOfWeek(y int64, m, d int) int {
	t := &time.Time{Year: y, Month: m, Day: d}
	t = time.SecondsToLocalTime(t.Seconds())
	return t.Weekday
}

// week in the year according to iso 8601
func calendarWeek(t *time.Time) int {
	dow := t.Weekday
	y, m, d := t.Year, t.Month, t.Day

	dow0101 := dayOfWeek(y, 1, 1)
	z := 0
	if dow0101 < 4 {
		z = 1
	}

	if m == 1 && 3 < dow0101 && dow0101 < 7-(d-1) {
		dow = dow0101 - 1
		dow0101 = dayOfWeek(y-1, 1, 1)
		m, d = 12, 31
	} else if m == 12 && 30-(d-1) < dayOfWeek(y+1, 1, 1) && dayOfWeek(y+1, 1, 1) < 4 {
		return 1
	}

	m--
	return z + 4*m + (2*m+(d-1)+dow0101-dow+6)*36/256
}

func FmtTime(sec int64, step TimeDelta) string {
	t := time.SecondsToLocalTime(sec)
	return step.Format(t)
}
