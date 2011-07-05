package chart

import (
	"fmt"
	"math"
	"time"
	//	"os"
	//	"strings"
)


type RangeExpansion int

const (
	ExpandNextTic = iota // Set min/max to next tic really below/above data-min/max data
	ExpandToTic          // Set to next tic below/above or equal to data-min/max data
	ExpandTight          // Use data min/max as limit 
	ExpandABit           // Like ExpandToTic and add/subtract half a tic distance.
)

type RangeMode struct {
	// If false: autoscaling. If true: use Value as fixed setting
	Fixed bool
	// If false: unconstrained autoscaling. If rue: use Lower and Upper as limits
	Constrained bool
	// One of ExpandNextTic, ExpandTight, ExpandABit
	Expand int
	// see above
	Value, Lower, Upper    float64
	TValue, TLower, TUpper *time.Time
}




type TicSetting struct {
	Hide   bool    // Dont show tics if true
	Minor  int     // 0: off, 1: clever, >1: number of intervalls
	Delta  float64 // Wanted step. 0 means auto 
	TDelta TimeDelta
	Fmt    string // special format string

}

type Tic struct {
	Pos, LabelPos float64
	Label         string
	Align         int  // -1: left, 0 center, 1 right (unused)
}


type Range struct {
	Log              bool      // logarithmic axis?
	Time             bool      // Time axis
	MinMode, MaxMode RangeMode // how to handel min and max
	TicSetting       TicSetting
	DataMin, DataMax float64 // actual values from data. if both zero: not calculated
	Min, Max         float64 // the min an d max of the xais
	TMin, TMax       *time.Time
	Tics             []Tic
	Norm             func(float64) float64 // map [Min:Max] to [0:1]
	InvNorm          func(float64) float64 // inverse of Norm()
	Data2Screen      func(float64) int
	Screen2Data      func(int) float64
}

var wochentage = []string{"So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"}

func calendarWeek(t *time.Time) int {
	// TODO(vodo): check if suitable
	jan01 := *t
	jan01.Month, jan01.Day, jan01.Hour, jan01.Minute, jan01.Second = 1, 1, 0, 0, 0
	diff := t.Seconds() - jan01.Seconds()
	week := int(float64(diff)/float64(60*60*24*7) + 0.5)
	if week == 0 {
		week++
	}
	return week
}

func FmtTime(sec int64, step TimeDelta) string {
	t := time.SecondsToLocalTime(sec)
	return step.Format(t)
}


var Units = []string{" y", " z", " a", " f", " p", " n", " µ", "m",
	" k", " M", " G", " T", " P", " E", " Z", " Y"}

func FmtFloat(f float64) string {
	af := math.Fabs(f)
	if f == 0 {
		return "0"
	} else if 0.1 <= af && af < 10 {
		return fmt.Sprintf("%.1f", f)
	} else if 10 <= af && af <= 1000 {
		return fmt.Sprintf("%.0f", f)
	}

	if af < 1 {
		var p = 8
		for math.Fabs(f) < 1 && p >= 0 {
			f *= 1000
			p--
		}
		return FmtFloat(f) + Units[p]
	} else {
		var p = 7
		for math.Fabs(f) > 1000 && p < 16 {
			f /= 1000
			p++
		}
		return FmtFloat(f) + Units[p]

	}
	return "xxx"
}


// Return val constrained by mode.
func Bound(mode RangeMode, val, ticDelta float64, upper bool) float64 {
	if mode.Fixed {
		return mode.Value
	}
	if mode.Constrained {
		if val < mode.Lower {
			val = mode.Lower
		} else if val > mode.Upper {
			val = mode.Upper
		}
	}
	switch mode.Expand {
	case ExpandToTic, ExpandNextTic:
		var v float64
		if upper {
			v = math.Ceil(val/ticDelta) * ticDelta
		} else {
			v = math.Floor(val/ticDelta) * ticDelta
		}
		if mode.Expand == ExpandNextTic && almostEqual(v, val) {
			if upper {
				v += ticDelta
			} else {
				v -= ticDelta
			}
		}
		val = v
	case ExpandABit:
		if upper {
			val += ticDelta / 2
		} else {
			val -= ticDelta / 2
		}
	}

	return val
}





// Return val constrained by mode.
func TimeBound(mode RangeMode, val *time.Time, step TimeDelta, upper bool) (bound *time.Time, tic *time.Time) {
	if mode.Fixed {
		bound = mode.TValue
		if upper {
			tic = RoundDown(val, step)
		} else {
			tic = RoundUp(val, step)
		}
		return
	}
	if mode.Constrained { // TODO(vodo) use T...
		sval := val.Seconds()
		if sval < int64(mode.Lower) {
			sval = int64(mode.Lower)
		} else if sval > int64(mode.Upper) {
			sval = int64(mode.Upper)
		}
		val = time.SecondsToLocalTime(sval)
	}

	switch mode.Expand {
	case ExpandToTic:
		if upper {
			val = RoundUp(val, step)
		} else {
			val = RoundDown(val, step)
		}
		return val, val
	case ExpandNextTic:
		if upper {
			tic = RoundUp(val, step)
		} else {
			tic = RoundDown(val, step)
		}
		s := tic.Seconds()
		if math.Fabs(float64(s-val.Seconds())/float64(step.Seconds())) < 0.15 {
			if upper {
				val = RoundUp(time.SecondsToLocalTime(s+step.Seconds()/2), step)
			} else {
				val = RoundDown(time.SecondsToLocalTime(s-step.Seconds()/2), step)
			}
		} else {
			val = tic
		}
		return val, val
	case ExpandABit:
		if upper {
			tic = RoundDown(val, step)
			val = time.SecondsToLocalTime(tic.Seconds() + step.Seconds()/2)
		} else {
			tic = RoundUp(val, step)
			val = time.SecondsToLocalTime(tic.Seconds() - step.Seconds()/2)
		}
		return

	}

	return val, val
}

func f2d(x float64) string {
	s := int64(x)
	t := time.SecondsToLocalTime(s)
	return t.Format("2006-01-02 15:04:05 (Mon)")
}



// Set up Range according to RangeModes and TicSettings.
// DataMin and DataMax should be present.
func (r *Range) Setup(desiredNumberOfTics, maxNumberOfTics, sWidth, sOffset int, revert bool) {
	if desiredNumberOfTics <= 1 {
		desiredNumberOfTics = 2
	}
	if maxNumberOfTics < desiredNumberOfTics {
		maxNumberOfTics = desiredNumberOfTics
	}
	if r.DataMax == r.DataMin {
		r.DataMax = r.DataMin + 1
	}
	delta := (r.DataMax - r.DataMin) / float64(desiredNumberOfTics-1)
	mindelta := (r.DataMax - r.DataMin) / float64(maxNumberOfTics-1)

	if r.Time {
		td := MatchingTimeDelta(delta, 3)
		mint := time.SecondsToLocalTime(int64(r.DataMin))
		maxt := time.SecondsToLocalTime(int64(r.DataMax))

		var ftic, ltic *time.Time
		r.TMin, ftic = TimeBound(r.MinMode, mint, td, false)
		r.TMax, ltic = TimeBound(r.MaxMode, maxt, td, true)
		r.TicSetting.Delta, r.TicSetting.TDelta = float64(td.Seconds()), td
		r.Min, r.Max = float64(r.TMin.Seconds()), float64(r.TMax.Seconds())

		ftd := float64(td.Seconds())
		actNumTics := int((r.Max - r.Min) / ftd)
		if actNumTics > maxNumberOfTics {
			fmt.Printf("Switching to next (%d > %d) delta from %s", actNumTics, maxNumberOfTics, td)
			td = NextTimeDelta(td)
			ftd = float64(td.Seconds())
			fmt.Printf("  -->  %s\n", td)
			r.TMin, ftic = TimeBound(r.MinMode, mint, td, false)
			r.TMax, ltic = TimeBound(r.MaxMode, maxt, td, true)
			r.TicSetting.Delta, r.TicSetting.TDelta = float64(td.Seconds()), td
			r.Min, r.Max = float64(r.TMin.Seconds()), float64(r.TMax.Seconds())
			actNumTics = int((r.Max - r.Min) / ftd)
			if actNumTics > maxNumberOfTics {
				fmt.Printf("Switching to next (%d > %d) delta from %s", actNumTics, maxNumberOfTics, td)
				td = NextTimeDelta(td)
				fmt.Printf("  -->  %s\n", td)
				r.TMin, ftic = TimeBound(r.MinMode, mint, td, false)
				r.TMax, ltic = TimeBound(r.MaxMode, maxt, td, true)
				r.TicSetting.Delta, r.TicSetting.TDelta = float64(td.Seconds()), td
				r.Min, r.Max = float64(r.TMin.Seconds()), float64(r.TMax.Seconds())
			}

		}

		fmt.Printf("Range:\n  Data:  %s  to  %s\n  --->   %s  to  %s\n  Tic-Delta: %s\n  Tics:  %s  to  %s\n",
			f2d(r.DataMin), f2d(r.DataMax), f2d(r.Min), f2d(r.Max), td,
			ftic.Format("2006-01-02 15:04:05 (Mon)"), ltic.Format("2006-01-02 15:04:05 (Mon)"))

		r.Tics = make([]Tic, 0)
		step := int64(td.Seconds())
		align := 0
		for i := 0; ftic.Seconds() <= ltic.Seconds(); i++ {
			x := float64(ftic.Seconds())
			label := td.Format(ftic)
			var labelPos float64
			if td.Period() {
				labelPos = x + float64(step)/2
			} else {
				labelPos = x
			}
			t := Tic{Pos: x, LabelPos: labelPos, Label: label, Align: align}
			r.Tics = append(r.Tics, t)
			// fmt.Printf("    Made Tic %s  '%s' %d at %s \n", t.Label, ftic.Format("2006-01-02 15:04:05 (Mon)"), align, time.SecondsToLocalTime(int64(t.LabelPos)).Format("2006-01-02 15:04:05 (Mon)") )
			ftic = RoundDown(time.SecondsToLocalTime(ftic.Seconds()+step+step/5), td)
			if i > maxNumberOfTics+3 {
				break
			}
		}

	} else { // simple, not a date range 
		de := math.Pow10(int(math.Floor(math.Log10(delta))))
		// fmt.Printf(":: %f, %f, %d \n", delta, de, int(de))
		f := delta / de
		switch {
		case f < 2:
			f = 1
		case f < 4:
			f = 2
		case f < 9:
			f = 5
		default:
			f = 1
			de *= 10
		}
		delta = f * de
		if delta < mindelta {
			switch f {
			case 1, 5:
				delta *= 2
			case 2:
				delta *= 2.5
			default:
				fmt.Printf("Oooops. Strange f: %g\n", f)
			}
		}

		r.Min = Bound(r.MinMode, r.DataMin, delta, false)
		r.Max = Bound(r.MaxMode, r.DataMax, delta, true)
		r.TicSetting.Delta = delta
		first := delta * math.Ceil(r.Min/delta)
		num := int(-first/delta + math.Floor(r.Max/delta) + 1.5)
		fmt.Printf("Range: (%g,%g) --> (%g,%g), Tic-Delta: %g, %d tics from %g\n", r.DataMin, r.DataMax, r.Min, r.Max, delta, num, first)

		r.Tics = make([]Tic, num)
		for i, x := 0, first; i < num; i, x = i+1, x+delta {
			r.Tics[i].Pos, r.Tics[i].LabelPos = x, x
			r.Tics[i].Label = FmtFloat(x)
		}

	}

	if !revert {
		r.Data2Screen = func(x float64) int {
			return int(math.Floor(float64(sWidth)*(x-r.Min)/(r.Max-r.Min))) + sOffset
		}
		r.Screen2Data = func(x int) float64 {
			return (r.Max-r.Min)*float64(x-sOffset)/float64(sWidth) + r.Min
		}
	} else {
		r.Data2Screen = func(x float64) int {
			return sWidth - int(math.Floor(float64(sWidth)*(x-r.Min)/(r.Max-r.Min))) + sOffset
		}
		r.Screen2Data = func(x int) float64 {
			return (r.Max-r.Min)*float64(-x+sOffset+sWidth)/float64(sWidth) + r.Min
		}

	}

}

type KeyEntry struct {
	Symbol int
	Text   string
}

type DataStyle struct {
	Symbol int     // -1: no symbol, 0: auto, 1... fixed
	Line   int     // 0: no line, 1, solid, 2 dashed, 3 dotted
	Size   float64 // 0: auto (1)
	Color  int     // index into palette
}


//        otl  otc  otr      
//       +--------------+ 
//   olt |itl  itc  itr | ort
//       |              |
//   olc |icl  icc  icr | ort
//       |              |
//   olb |ibl  ibc  ibr | orb
//       +--------------+ 
//        obl  obc  obr
type Key struct {
	Hide   bool
	Layout struct {
		Cols, Rows int // 0,0 means 1,1
	}
	Border int    // -1: off, 0: std, 1...:other styles
	Pos    string // "": itr
	// Width, Height int    // 0,0: auto
	X, Y int
}

type ChartValue interface {
	// center
	CX() float64
	CY() float64
	// bounding box
	MinX() float64
	MinY() float64
	MaxX() float64
	MaxY() float64
}

// Simple Point
type Point struct{ X, Y float64 }

func (p Point) CX() float64   { return p.X }
func (p Point) CY() float64   { return p.Y }
func (p Point) MinX() float64 { return p.X }
func (p Point) MinY() float64 { return p.Y }
func (p Point) MaxX() float64 { return p.X }
func (p Point) MaxY() float64 { return p.Y }

// Box in Boxplot
type Box struct {
	X, Avg, Med, Q1, Q3, Low, High float64
	Outliers                       []float64
}

func (p Box) CX() float64   { return p.X }
func (p Box) CY() float64   { return p.Med }
func (p Box) MinX() float64 { return p.X }
func (p Box) MinY() float64 {
	x := minimum(p.Outliers)
	if x != math.NaN() {
		return x
	}
	return p.Low
}
func (p Box) MaxX() float64 { return p.X }
func (p Box) MaxY() float64 {
	x := maximum(p.Outliers)
	if x != math.NaN() {
		return x
	}
	return p.High
}


/***** bad idea *******
// Bin in Histograms
type Bin struct {
	X, Width float64
	Count    int
}

func (p Bin) CX() float64   { return p.X }
func (p Bin) CY() float64   { return float64(p.Count) / 2 }
func (p Bin) MinX() float64 { return p.X - p.Width/2 }
func (p Bin) MinY() float64 { return 0 }
func (p Bin) MaxX() float64 { return p.X + p.Width/2 }
func (p Bin) MaxY() float64 { return float64(p.Count) }
***********************/

type ChartData struct {
	Name   string
	Style  DataStyle
	Values []ChartValue
}
