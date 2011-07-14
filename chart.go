package chart

import (
	"fmt"
	"math"
	"time"
	//	"os"
	"strings"
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

	Min, Max   float64    // Minium and Maximum of this axis/range.
	TMin, TMax *time.Time // Same as Min/Max, but used for Date/Time axis

	Norm        func(float64) float64 // Function to map [Min:Max] to [0:1]
	InvNorm     func(float64) float64 // Inverse of Norm()
	Data2Screen func(float64) int     // Function to map data value to screen position
	Screen2Data func(int) float64     // Inverse of Data2Screen
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

	fmt.Printf("Range:\n  Data:  %s  to  %s\n  --->   %s  to  %s\n  Tic-Delta: %s\n  Tics:  %s  to  %s\n",
		f2d(r.DataMin), f2d(r.DataMax), f2d(r.Min), f2d(r.Max), td,
		ftic.Format("2006-01-02 15:04:05 (Mon)"), ltic.Format("2006-01-02 15:04:05 (Mon)"))

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
		fmt.Printf("Range: (%g,%g) --> (%g,%g), Tic-Delta: %g, %d tics from %g\n",
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
//   Norm and InvNorm     TODO(vodo) not jet implemented
//   Data2Screen
//   Screen2Data
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

var Palette = []string{"#000000", // 0: black
	"#dd0000", "#0000dd", "#00bb00", // 1: red, 2: green, 3: blue
	"#bb00bb", "#00aaaa", "#aaaa00", // 4: purple, 5: turkis, 6: gold
	"#ff8888", "#8888ff", "#66ff66", // 7: light red, 8: light blue, 9: light green
	// Last are system colors
	"#000000",                                                                   // len-9: black
	"#202020", "#404040", "#606060", "#808080", "#a0a0a0", "#c0c0c0", "#e0e0e0", // len-2 to len-8: gray
	"#ffffff", // len-1: white
}

// DataStyle contains all information about all graphic elements in a chart.
type DataStyle struct {
	Symbol      int     // -1: no symbol, 0: auto, 1... fixed
	SymbolColor int     // 0: auto
	Line        int     // -1: no line, 0: auto, 1, solid, 2 dashed, 3 dotted, 4 dashdotdot, 5 longdash 6 longdot
	LineColor   int     // 0: auto = same as SymbolColor
	Size        float64 // 0: auto (=1)
	Font        string  // "": default
	FontSize    int     // -2: tiny, -1: small, 0: normal, 1: large, 2: huge
	Alpha       float64
}

// Next yields the next style: Next symbol, line type and colors.
func (ds DataStyle) Next() DataStyle {
	ds.Symbol = (ds.Symbol + 1) % len(Symbol)
	ds.SymbolColor = (ds.SymbolColor + 1) % (len(Palette) - 9) // 9 system colors
	ds.Line = (ds.Line+1)%6 + 1
	if ds.LineColor != 0 {
		ds.LineColor = (ds.LineColor + 1) % (len(Palette) - 9) // 9 system colors
	}
	return ds
}

// NextMerge yields the next data style and merges all "non-auto" setting s of.
func (ds DataStyle) NextMerge(m DataStyle) DataStyle {
	ds = ds.Next()
	if m.Symbol != 0 {
		ds.Symbol = m.Symbol
	}
	if m.SymbolColor != 0 {
		ds.SymbolColor = m.SymbolColor
	}
	if m.Line != 0 {
		ds.Line = m.Line
	}
	if m.LineColor != 0 {
		ds.LineColor = m.LineColor
	}
	if m.Size != 0 {
		ds.Size = m.Size
	}
	if m.Font != "" {
		ds.Font = m.Font
	}
	if m.FontSize != 0 {
		ds.FontSize = m.FontSize
	}
	if m.Alpha != 0 {
		ds.Alpha = m.Alpha
	}
	return ds
}


// Key encapsulates settings for keys/legends in a chart.
//
// Key placement os governed by Pos which may take the following values:
//          otl  otc  otr      
//         +--------------+ 
//     olt |itl  itc  itr | ort
//         |              |
//     olc |icl  icc  icr | ort
//         |              |
//     olb |ibl  ibc  ibr | orb
//         +--------------+ 
//        obl  obc  obr
//
type Key struct {
	Hide    bool       // Don't show key/legend if true
	Cols    int        // Number of colums to use. If <0 fill rows before colums
	Border  int        // -1: off, 0: std, 1...:other styles
	Pos     string     // default "" is "itr"
	X, Y    int        // Coordiantes where to put in chart.
	Entries []KeyEntry // List of entries in the legend
}


// KeyEntry encapsulates an antry in the key/legend.
type KeyEntry struct {
	Symbol int    // Symbol index to use
	Linie  int    // Line Style
	Text   string // Text to display
}

// Margins
var KL_LRBorder int = 1 // before and after whole key
var KL_SLSep int = 2    // space between symbol and test
var KL_ColSep int = 2   // space between columns
var KL_MLSep int = 1    // extra space between rows if multiline text are present

func (key *Key) LayoutKeyTxt() (kb *TextBuf) {
	// TODO(vodo) the following is ugly (and stinks)
	if key.Hide {
		return
	}

	// count real entries in num, see if multilines are present in haveml
	num, haveml := 0, false
	for _, e := range key.Entries {
		if e.Text == "" {
			continue
		}
		num++
		lines := strings.Split(e.Text, "\n", -1)
		if len(lines) > 1 {
			haveml = true
		}
	}
	if num == 0 {
		return
	} // no entries

	rowfirst := false
	cols := key.Cols
	if cols < 0 {
		cols = -cols
		rowfirst = true
	}
	if cols == 0 {
		cols = 1
	}
	if num < cols {
		cols = num
	}
	rows := (num + cols - 1) / cols

	// fmt.Printf("%d entries on %d columns: %d rows\n", num, cols, rows)

	// Arrays with infos
	width := make([][]int, cols)
	for i := 0; i < cols; i++ {
		width[i] = make([]int, rows)
	}
	height := make([][]int, cols)
	for i := 0; i < cols; i++ {
		height[i] = make([]int, rows)
	}
	symbol := make([][]int, cols)
	for i := 0; i < cols; i++ {
		symbol[i] = make([]int, rows)
	}
	text := make([][][]string, cols)
	for i := 0; i < cols; i++ {
		text[i] = make([][]string, rows)
	}

	// fill arrays
	i := 0
	for _, e := range key.Entries {
		if e.Text == "" {
			continue
		}
		var r, c int
		if rowfirst {
			r, c = i/cols, i%cols
		} else {
			c, r = i/rows, i%rows
		}
		lines := strings.Split(e.Text, "\n", -1)
		ml := 0
		for _, t := range lines {
			if len(t) > ml { // TODO(vodo) use utf8.CountRuneInString and honour different chars
				ml = len(t)
			}
		}
		symbol[c][r] = e.Symbol // TODO(vodo) allow line symbols?
		height[c][r] = len(lines)
		width[c][r] = ml
		text[c][r] = lines
		i++
	}
	colwidth := make([]int, cols)
	rowheight := make([]int, rows)
	totalheight, totalwidth := 0, 0
	for c := 0; c < cols; c++ {
		max := 0
		for r := 0; r < rows; r++ {
			if width[c][r] > max {
				max = width[c][r]
			}
		}
		max += 2*KL_LRBorder + 1 + KL_SLSep // formt is " *  Label "
		colwidth[c] = max
		totalwidth += max
	}
	for r := 0; r < rows; r++ {
		max := 0
		for c := 0; c < cols; c++ {
			if height[c][r] > max {
				max = height[c][r]
			}
		}
		rowheight[r] = max
		totalheight += max
	}

	// width and height: + 2 for outer border/box
	w := totalwidth + KL_ColSep*(cols-1) + 2
	h := totalheight + 2
	if haveml {
		h += KL_MLSep * (rows - 1)
	}
	kb = NewTextBuf(w, h)
	if key.Border != -1 {
		kb.Rect(0, 0, w-1, h-1, key.Border+1, ' ')
	}

	// Produce box
	x := 1
	for c := 0; c < cols; c++ {
		y := 1
		for r := 0; r < rows; r++ {
			if width[c][r] == 0 {
				continue
			}
			xx := x + KL_LRBorder
			if symbol[c][r] != -1 {
				kb.Put(xx, y, symbol[c][r])
				xx += 1 + KL_SLSep
			}
			for l, t := range text[c][r] {
				kb.Text(xx, y+l, t, -1)
			}
			y += rowheight[r]
			if haveml {
				y += KL_MLSep
			}
		}
		x += colwidth[c] + KL_ColSep
	}

	return
}


func LayoutTxt(w, h int, title, xlabel, ylabel string, hidextics, hideytics bool, key *Key) (width, leftm, height, topm int, kb *TextBuf, numxtics, numytics int) {
	if key.Pos == "" {
		key.Pos = "itr"
	}

	if h < 5 {
		h = 5
	}
	if w < 10 {
		w = 10
	}

	width, leftm, height, topm = w-6, 2, h-1, 0
	xlabsep, ylabsep := 1, 3
	if title != "" {
		topm++
		height--
	}
	if xlabel != "" {
		height--
	}
	if !hidextics {
		height--
		xlabsep++
	}
	if ylabel != "" {
		leftm += 2
		width -= 2
	}
	if !hideytics {
		leftm += 6
		width -= 6
		ylabsep += 6
	}

	if !key.Hide {
		kb = key.LayoutKeyTxt()
		if kb != nil {
			kw, kh := kb.W, kb.H
			switch key.Pos[:2] {
			case "ol":
				width, leftm = width-kw-2, leftm+kw
				key.X = 0
			case "or":
				width = width - kw - 2
				key.X = w - kw
			case "ot":
				height, topm = height-kh-2, topm+kh
				key.Y = 1
			case "ob":
				height = height - kh - 2
				key.Y = topm + height + 4
			case "it":
				key.Y = topm + 1
			case "ic":
				key.Y = topm + (height-kh)/2
			case "ib":
				key.Y = topm + height - kh

			}

			switch key.Pos[:2] {
			case "ol", "or":
				switch key.Pos[2] {
				case 't':
					key.Y = topm
				case 'c':
					key.Y = topm + (height-kh)/2
				case 'b':
					key.Y = topm + height - kh + 1
				}
			case "ot", "ob":
				switch key.Pos[2] {
				case 'l':
					key.X = leftm
				case 'c':
					key.X = leftm + (width-kw)/2
				case 'r':
					key.X = w - kw - 2
				}
			}
			if key.Pos[0] == 'i' {
				switch key.Pos[2] {
				case 'l':
					key.X = leftm + 2
				case 'c':
					key.X = leftm + (width-kw)/2
				case 'r':
					key.X = leftm + width - kw - 2
				}

			}
		}
	}

	// fmt.Printf("width=%d, height=%d, leftm=%d, topm=%d\n", width, height, leftm, topm)

	switch {
	case width < 20:
		numxtics = 2
	case width < 30:
		numxtics = 3
	case width < 60:
		numxtics = 4
	case width < 80:
		numxtics = 5
	case width < 100:
		numxtics = 7
	default:
		numxtics = 10
	}
	// fmt.Printf("Requesting %d,%d tics.\n", ntics,height/3)

	numytics = h / 5

	return
}


// Print xrange to tb at vertical position y.
// Axis, tics, tic labels, axis label and range limits are drawn.
// mirror: 0: no other axis, 1: axis without tics, 2: axis with tics,
func TxtXRange(xrange Range, tb *TextBuf, y, y1 int, label string, mirror int) {
	xa, xe := xrange.Data2Screen(xrange.Min), xrange.Data2Screen(xrange.Max)
	for sx := xa; sx <= xe; sx++ {
		tb.Put(sx, y, '-')
		if mirror >= 1 {
			tb.Put(sx, y1, '-')
		}
	}
	if xrange.ShowZero && xrange.Min < 0 && xrange.Max > 0 {
		z := xrange.Data2Screen(0)
		for yy := y - 1; yy > y1+1; yy-- {
			tb.Put(z, yy, ':')
		}
	}

	if label != "" {
		yy := y + 1
		if !xrange.TicSetting.Hide {
			yy++
		}
		tb.Text((xa+xe)/2, yy, label, 0)
	}

	for _, tic := range xrange.Tics {
		x := xrange.Data2Screen(tic.Pos)
		lx := xrange.Data2Screen(tic.LabelPos)
		if xrange.Time {
			tb.Put(x, y, '|')
			if mirror >= 2 {
				tb.Put(x, y1, '|')
			}
			tb.Put(x, y+1, '|')
			if tic.Align == -1 {
				tb.Text(lx+1, y+1, tic.Label, -1)
			} else {
				tb.Text(lx, y+1, tic.Label, 0)
			}
		} else {
			tb.Put(x, y, '+')
			if mirror >= 2 {
				tb.Put(x, y1, '+')
			}
			tb.Text(lx, y+1, tic.Label, 0)
		}
		if xrange.ShowLimits {
			if xrange.Time {
				tb.Text(xa, y+2, xrange.TMin.Format("2006-01-02 15:04:05"), -1)
				tb.Text(xe, y+2, xrange.TMax.Format("2006-01-02 15:04:05"), 1)
			} else {
				tb.Text(xa, y+2, fmt.Sprintf("%g", xrange.Min), -1)
				tb.Text(xe, y+2, fmt.Sprintf("%g", xrange.Max), 1)
			}
		}
	}
}


// Print yrange to tb at horizontal position x.
// Axis, tics, tic labels, axis label and range limits are drawn.
// mirror: 0: no other axis, 1: axis without tics, 2: axis with tics,
func TxtYRange(yrange Range, tb *TextBuf, x, x1 int, label string, mirror int) {
	ya, ye := yrange.Data2Screen(yrange.Min), yrange.Data2Screen(yrange.Max)
	for sy := min(ya, ye); sy <= max(ya, ye); sy++ {
		tb.Put(x, sy, '|')
		if mirror >= 1 {
			tb.Put(x1, sy, '|')
		}
	}
	if yrange.ShowZero && yrange.Min < 0 && yrange.Max > 0 {
		z := yrange.Data2Screen(0)
		for xx := x + 1; xx < x1; xx += 2 {
			tb.Put(xx, z, '-')
		}
	}

	if label != "" {
		tb.Text(1, (ya+ye)/2, label, 3)
	}

	for _, tic := range yrange.Tics {
		y := yrange.Data2Screen(tic.Pos)
		ly := yrange.Data2Screen(tic.LabelPos)
		if yrange.Time {
			tb.Put(x, y, '+')
			if mirror >= 2 {
				tb.Put(x1, y, '+')
			}
			if tic.Align == 0 { // centered tic
				tb.Put(x-1, y, '-')
				tb.Put(x-2, y, '-')
			}
			tb.Text(x, ly, tic.Label+" ", 1)
		} else {
			tb.Put(x, y, '+')
			if mirror >= 2 {
				tb.Put(x1, y, '+')
			}
			tb.Text(x-2, ly, tic.Label, 1)
		}
	}
}
