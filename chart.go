package chart

import (
	"fmt"
	"math"
	"time"
	//	"os"
)


// Suitable values for Expand in RangeMode.
const (
	ExpandNextTic = iota // Set min/max to next tic really below/above min/max of data
	ExpandToTic          // Set to next tic below/above or equal to min/max of data
	ExpandTight          // Use data min/max as limit 
	ExpandABit           // Like ExpandToTic and add/subtract ExpandABitFraction of tic distance.
)

var ExpandABitFraction = 0.5 // Fraction of tic spacing added in ExpandABit Range.Expand mode.

// RangeMode describes how one end of an axis is set up. There are basically three different main modes:
//   o Fixed: Fixed==true. 
//     Use Value/TValue as fixed value ignoring data 
//   o Unconstrained autoscaling: Fixed==false && Constrained==false
//     Set range to whatever data requires
//   o Constrained autoscaling: Fixed==false && Constrained==true
//     Scale axis according to data present, but limit scaling to intervall [Lower,Upper]
// For both autoscaling modes Expand defines how much expansion is done below/above
// the lowest/highest data point.
type RangeMode struct {
	Fixed          bool       // If false: autoscaling. If true: use (T)Value/TValue as fixed setting
	Constrained    bool       // If false: full autoscaling. If true: use (T)Lower (T)Upper as limits
	Expand         int        // One of ExpandNextTic, ExpandTight, ExpandABit
	Value          float64    // Value of end point of axis in Fixed=true mode, ignorder otherwise
	TValue         *time.Time // Same as Value, but used for Date/Time axis
	Lower, Upper   float64    // Lower and upper limit for constrained autoscaling
	TLower, TUpper *time.Time // Same s Lower/Upper, but used for Date/Time axis
}


// TicSettings describes how (if at all) tics are shown on an axis.
type TicSetting struct {
	Hide   bool      // Dont show tics if true
	Minor  int       // 0: off, 1: auto, >1: number of intervalls (not number of tics!)
	Delta  float64   // Wanted step. 0 means auto 
	TDelta TimeDelta // Same as Delta, used for Date/Time axis
	Fmt    string    // special format string
	Grid   int       // 0: none, 1: lines, 2: blocks
	Mirror int       // 0: mirror axis and tics, -1: don't mirror anything, 1: mirror axis only (no tics)
}

// Tic describs a single tic on an axis.
type Tic struct {
	Pos, LabelPos float64 // Position if the tic and its label
	Label         string  // The Label
	Align         int     // Alignment of the label: -1: left/top, 0 center, 1 right/bottom (unused)
}

// Range encapsulates all information about an axis.
type Range struct {
	Log              bool       // logarithmic axis?
	Time             bool       // Date/Time axis
	MinMode, MaxMode RangeMode  // How to handel min and max of this axis/range
	TicSetting       TicSetting // How to handle tics.
	DataMin, DataMax float64    // Actual min/max values from data. If both zero: not calculated
	ShowLimits       bool       // Display axis Min and Max on plot
	ShowZero         bool       // Add line to show 0 of this axis
	Tics             []Tic      // List of tics to display
	Label            string     // Label of axis

	Min, Max   float64    // Minium and Maximum of this axis/range.
	TMin, TMax *time.Time // Same as Min/Max, but used for Date/Time axis

	Norm        func(float64) float64 // Function to map [Min:Max] to [0:1]
	InvNorm     func(float64) float64 // Inverse of Norm()
	Data2Screen func(float64) int     // Function to map data value to screen position
	Screen2Data func(int) float64     // Inverse of Data2Screen
}

// Prepare the range r for use, especially set up all values needed for autoscale() to work properly
func (r *Range) init() {
	// All the min stuff
	if r.MinMode.Fixed {
		// copy TValue to Value if set and time axis
		if r.Time && r.MinMode.TValue != nil {
			r.MinMode.Value = float64(r.MinMode.TValue.Seconds())
		}
		r.DataMin = r.MinMode.Value
	} else if r.MinMode.Constrained {
		// copy TLower/TUpper to Lower/Upper if set and time axis
		if r.Time && r.MinMode.TLower != nil {
			r.MinMode.Lower = float64(r.MinMode.TLower.Seconds())
		}
		if r.Time && r.MinMode.TUpper != nil {
			r.MinMode.Upper = float64(r.MinMode.TUpper.Seconds())
		}
		if r.MinMode.Lower == 0 && r.MinMode.Upper == 0 {
			// Constrained but un-initialized: Full autoscaling
			r.MinMode.Lower = -math.MaxFloat64
			r.MinMode.Upper = math.MaxFloat64
		}
		r.DataMin = r.MinMode.Upper
	} else {
		r.DataMin = math.MaxFloat64
	}

	// All the max stuff
	if r.MaxMode.Fixed {
		// copy TValue to Value if set and time axis
		if r.Time && r.MaxMode.TValue != nil {
			r.MaxMode.Value = float64(r.MaxMode.TValue.Seconds())
		}
		r.DataMax = r.MaxMode.Value
	} else if r.MaxMode.Constrained {
		// copy TLower/TUpper to Lower/Upper if set and time axis
		if r.Time && r.MaxMode.TLower != nil {
			r.MaxMode.Lower = float64(r.MaxMode.TLower.Seconds())
		}
		if r.Time && r.MaxMode.TUpper != nil {
			r.MaxMode.Upper = float64(r.MaxMode.TUpper.Seconds())
		}
		if r.MaxMode.Lower == 0 && r.MaxMode.Upper == 0 {
			// Constrained but un-initialized: Full autoscaling
			r.MaxMode.Lower = -math.MaxFloat64
			r.MaxMode.Upper = math.MaxFloat64
		}
		r.DataMax = r.MaxMode.Upper
	} else {
		r.DataMax = -math.MaxFloat64
	}

	fmt.Printf("At end of init: DataMin / DataMax  =   %g / %g\n", r.DataMin, r.DataMax)
}


// Update DataMin and DataMax according to the RangeModes.
func (r *Range) autoscale(x float64) {

	if x < r.DataMin && !r.MinMode.Fixed {
		if !r.MinMode.Constrained {
			// full autoscaling
			r.DataMin = x
		} else {
			r.DataMin = fmin(fmax(x, r.MinMode.Lower), r.DataMin)
		}
	}

	if x > r.DataMax && !r.MaxMode.Fixed {
		if !r.MaxMode.Constrained {
			// full autoscaling
			r.DataMax = x
		} else {
			r.DataMax = fmax(fmin(x, r.MaxMode.Upper), r.DataMax)
		}
	}
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


var Units = []string{" y", " z", " a", " f", " p", " n", " µ", "m", " k", " M", " G", " T", " P", " E", " Z", " Y"}

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

func almostEqual(a, b, d float64) bool {
	return math.Fabs(a-b) < d
}


// ApplyRangeMode returns val constrained by mode. val is considered the upper end of an range/axis
// if upper is true. To allow proper rounding to tic (depending on desired RangeMode)
// the ticDelta has to be provided. Logaritmic axis are selected by log = true and ticDelta
// is ignored: Tics are of the form 1*10^n.
func ApplyRangeMode(mode RangeMode, val, ticDelta float64, upper, log bool) float64 {
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
			if log {
				v = math.Pow10(int(math.Ceil(math.Log10(val))))
			} else {
				v = math.Ceil(val/ticDelta) * ticDelta
			}
		} else {
			if log {
				v = math.Pow10(int(math.Floor(math.Log10(val))))
			} else {
				v = math.Floor(val/ticDelta) * ticDelta
			}
		}
		if mode.Expand == ExpandNextTic {
			if upper {
				if log {
					if val/v < 2 { // TODO(vodo) use ExpandABitFraction
						v *= ticDelta
					}
				} else {
					if almostEqual(v, val, ticDelta/15) {
						v += ticDelta
					}
				}
			} else {
				if log {
					if v/val > 7 { // TODO(vodo) use ExpandABitFraction
						v /= ticDelta
					}
				} else {
					if almostEqual(v, val, ticDelta/15) {
						v -= ticDelta
					}
				}
			}
		}
		val = v
	case ExpandABit:
		if upper {
			if log {
				val *= math.Pow(10, ExpandABitFraction)
			} else {
				val += ticDelta * ExpandABitFraction
			}
		} else {
			if log {
				val /= math.Pow(10, ExpandABitFraction)
			} else {
				val -= ticDelta * ExpandABitFraction
			}
		}
	}

	return val
}


// TApplyRangeMode is the same as ApplyRangeMode for date/time axis/ranges.
func TApplyRangeMode(mode RangeMode, val *time.Time, step TimeDelta, upper bool) (bound *time.Time, tic *time.Time) {
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


func (r *Range) tSetup(desiredNumberOfTics, maxNumberOfTics int, delta, mindelta float64) {
	var td TimeDelta
	if r.TicSetting.TDelta != nil {
		td = r.TicSetting.TDelta
	} else {
		td = MatchingTimeDelta(delta, 3)
	}
	r.ShowLimits = true

	// Set up time tic delta
	mint := time.SecondsToLocalTime(int64(r.DataMin))
	maxt := time.SecondsToLocalTime(int64(r.DataMax))

	var ftic, ltic *time.Time
	r.TMin, ftic = TApplyRangeMode(r.MinMode, mint, td, false)
	r.TMax, ltic = TApplyRangeMode(r.MaxMode, maxt, td, true)
	r.TicSetting.Delta, r.TicSetting.TDelta = float64(td.Seconds()), td
	r.Min, r.Max = float64(r.TMin.Seconds()), float64(r.TMax.Seconds())

	ftd := float64(td.Seconds())
	actNumTics := int((r.Max - r.Min) / ftd)
	if actNumTics > maxNumberOfTics {
		// recalculate time tic delta
		fmt.Printf("Switching to next (%d > %d) delta from %s", actNumTics, maxNumberOfTics, td)
		td = NextTimeDelta(td)
		ftd = float64(td.Seconds())
		fmt.Printf("  -->  %s\n", td)
		r.TMin, ftic = TApplyRangeMode(r.MinMode, mint, td, false)
		r.TMax, ltic = TApplyRangeMode(r.MaxMode, maxt, td, true)
		r.TicSetting.Delta, r.TicSetting.TDelta = float64(td.Seconds()), td
		r.Min, r.Max = float64(r.TMin.Seconds()), float64(r.TMax.Seconds())
		actNumTics = int((r.Max - r.Min) / ftd)
	}

	/*
		fmt.Printf("Range:\n  Data:  %s  to  %s\n  --->   %s  to  %s\n  Tic-Delta: %s\n  Tics:  %s  to  %s\n",
			f2d(r.DataMin), f2d(r.DataMax), f2d(r.Min), f2d(r.Max), td,
			ftic.Format("2006-01-02 15:04:05 (Mon)"), ltic.Format("2006-01-02 15:04:05 (Mon)"))
	*/
	// Set up tics
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
		ftic = RoundDown(time.SecondsToLocalTime(ftic.Seconds()+step+step/5), td)
		if i > maxNumberOfTics+3 {
			break
		}
	}
}

func (r *Range) fDelta(delta, mindelta float64) float64 {
	fmt.Printf("fDelta(%.3f, %.3f)\n", delta, mindelta)
	if r.Log {
		return 10
	}

	// Set up nice tic delta of the form 1,2,5 * 10^n
	de := math.Pow10(int(math.Floor(math.Log10(delta))))
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
		fmt.Printf("Redoing delta")
		// recalculate tic delta
		switch f {
		case 1, 5:
			delta *= 2
		case 2:
			delta *= 2.5
		default:
			fmt.Printf("Oooops. Strange f: %g\n", f)
		}
	}
	return delta
}

func (r *Range) fSetup(desiredNumberOfTics, maxNumberOfTics int, delta, mindelta float64) {
	if r.TicSetting.Delta != 0 {
		delta = r.TicSetting.Delta
	} else {
		delta = r.fDelta(delta, mindelta)
	}

	r.Min = ApplyRangeMode(r.MinMode, r.DataMin, delta, false, r.Log)
	r.Max = ApplyRangeMode(r.MaxMode, r.DataMax, delta, true, r.Log)
	r.TicSetting.Delta = delta
	if r.Log {
		x := math.Pow10(int(math.Ceil(math.Log10(r.Min))))
		last := math.Pow10(int(math.Floor(math.Log10(r.Max))))
		r.Tics = make([]Tic, 0, maxNumberOfTics)
		for ; x <= last; x = x * delta {
			t := Tic{Pos: x, LabelPos: x, Label: FmtFloat(x)}
			r.Tics = append(r.Tics, t)
			// fmt.Printf("%v\n", t)
		}

	} else {
		first := delta * math.Ceil(r.Min/delta)
		num := int(-first/delta + math.Floor(r.Max/delta) + 1.5)
		fmt.Printf("Range: (%.2f,%.2f) --> (%g,%g), Tic-Delta: %g, %d tics from %g\n",
			r.DataMin, r.DataMax, r.Min, r.Max, delta, num, first)

		// Set up tics
		r.Tics = make([]Tic, num)
		for i, x := 0, first; i < num; i, x = i+1, x+delta {
			r.Tics[i].Pos, r.Tics[i].LabelPos = x, x
			r.Tics[i].Label = FmtFloat(x)
		}

		// TODO(vodo) r.ShowLimits = true
	}
}


// SetUp sets up several fields of Range r according to RangeModes and TicSettings.
// DataMin and DataMax of r must be present and should indicate lowest and highest
// value present in the data set. The following field if r are filled:
//   (T)Min and (T)Max    lower and upper limit of axis, (T)-version for date/time axis
//   Tics                 slice of tics to draw
//   TicSetting.(T)Delta  actual tic delta
//   Norm and InvNorm     mapping of [lower,upper]_data --> [0:1] and inverse
//   Data2Screen          mapping of data to screen coordinates
//   Screen2Data          inverse of Data2Screen
// The parameters desired- and maxNumberOfTics are what the say.
// sWidth and sOffset are screen-width and -offset and are used to set up the
// Data-Screen conversion functions. If revert is true, than screen coordinates
// are asumed to be the other way around than mathematical coordinates.
//
// TODO(vodo) seperate screen stuff into own method.
func (r *Range) Setup(desiredNumberOfTics, maxNumberOfTics, sWidth, sOffset int, revert bool) {
	// Sanitize input
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

	fmt.Printf("Data: [%.2f:%.2f] --> delta/mindelta = %.2f/%.2f (desired %d/max %d)\n",
		r.DataMin, r.DataMax, delta, mindelta, desiredNumberOfTics, maxNumberOfTics)

	if r.Time {
		r.tSetup(desiredNumberOfTics, maxNumberOfTics, delta, mindelta)
	} else { // simple, not a date range 
		r.fSetup(desiredNumberOfTics, maxNumberOfTics, delta, mindelta)
	}

	if r.Log {
		r.Norm = func(x float64) float64 { return math.Log10(x/r.Min) / math.Log10(r.Max/r.Min) }
		r.InvNorm = func(f float64) float64 { return (r.Max-r.Min)*f + r.Min }
	} else {
		r.Norm = func(x float64) float64 { return (x - r.Min) / (r.Max - r.Min) }
		r.InvNorm = func(f float64) float64 { return (r.Max-r.Min)*f + r.Min }
	}

	if !revert {
		r.Data2Screen = func(x float64) int {
			return int(float64(sWidth)*r.Norm(x)) + sOffset
		}
		r.Screen2Data = func(x int) float64 {
			return r.InvNorm(float64(x-sOffset) / float64(sWidth))
		}
	} else {
		r.Data2Screen = func(x float64) int {
			return sWidth - int(float64(sWidth)*r.Norm(x)) + sOffset
		}
		r.Screen2Data = func(x int) float64 {
			return r.InvNorm(float64(-x+sOffset+sWidth) / float64(sWidth))
		}

	}

}


// LayoutData encapsulates the layout of the graph area in the whole drawing area.
type LayoutData struct {
	Width, Height      int // width and height of graph area
	Left, Top          int // left and top margin
	KeyX, KeyY         int // x and y coordiante of key
	NumXtics, NumYtics int // suggested numer of tics for both axis
}


// TODO: Key.X/Y have to go to explicit data
func Layout(g Graphics, title, xlabel, ylabel string, hidextics, hideytics bool, key *Key) (ld LayoutData) {
	fw, fh, _ := g.FontMetrics(g.Font("key"))
	w, h := g.Dimensions()

	if key.Pos == "" {
		key.Pos = "itr"
	}

	width, leftm, height, topm := w-int(6*fw), int(2*fw), h-2*fh, fh
	xlabsep, ylabsep := fh, int(3*fw)
	if title != "" {
		topm += (5 * fh) / 2
		height -= (5 * fh) / 2
	}
	if xlabel != "" {
		height -= (3 * fh) / 2
	}
	if !hidextics {
		height -= (3 * fh) / 2
		xlabsep += (3 * fh) / 2
	}
	if ylabel != "" {
		leftm += 2 * fh
		width -= 2 * fh
	}
	if !hideytics {
		leftm += int(6 * fw)
		width -= int(6 * fw)
		ylabsep += int(6 * fw)
	}

	if key != nil && !key.Hide && len(key.Place()) > 0 {
		m := key.Place()
		kw, kh, _, _ := key.Layout(g, m)
		sepx, sepy := int(fw)+fh, int(fw)+fh
		switch key.Pos[:2] {
		case "ol":
			width, leftm = width-kw-sepx, leftm+kw
			ld.KeyX = sepx / 2
		case "or":
			width = width - kw - sepx
			ld.KeyX = w - kw - sepx/2
		case "ot":
			height, topm = height-kh-sepy, topm+kh
			ld.KeyY = sepy / 2
		case "ob":
			height = height - kh - sepy
			ld.KeyY = h - kh - sepy/2
		case "it":
			ld.KeyY = topm + sepy
		case "ic":
			ld.KeyY = topm + (height-kh)/2
		case "ib":
			ld.KeyY = topm + height - kh - sepy

		}

		switch key.Pos[:2] {
		case "ol", "or":
			switch key.Pos[2] {
			case 't':
				ld.KeyY = topm
			case 'c':
				ld.KeyY = topm + (height-kh)/2
			case 'b':
				ld.KeyY = topm + height - kh
			}
		case "ot", "ob":
			switch key.Pos[2] {
			case 'l':
				ld.KeyX = leftm
			case 'c':
				ld.KeyX = leftm + (width-kw)/2
			case 'r':
				ld.KeyX = w - kw - sepx
			}
		}
		if key.Pos[0] == 'i' {
			switch key.Pos[2] {
			case 'l':
				ld.KeyX = leftm + sepx
			case 'c':
				ld.KeyX = leftm + (width-kw)/2
			case 'r':
				ld.KeyX = leftm + width - kw - sepx
			}
		}
	}

	// fmt.Printf("width=%d, height=%d, leftm=%d, topm=%d  (fw=%d)\n", width, height, leftm, topm, int(fw))

	// Number of tics
	if width/int(fw) <= 20 {
		ld.NumXtics = 2
	} else {
		ld.NumXtics = width / int(8*fw)
		if ld.NumXtics > 25 {
			ld.NumXtics = 25
		}
	}
	ld.NumYtics = height / (4 * fh)
	if ld.NumYtics > 20 {
		ld.NumYtics = 20
	}

	ld.Width, ld.Height = width, height
	ld.Left, ld.Top = leftm, topm

	return
}
