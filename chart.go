package chart

import (
	"fmt"
	"math"
	"time"
	//	"os"
	"strings"
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
	Align         int // -1: left/top, 0 center, 1 right/bottom (unused)
}


type Range struct {
	Log              bool      // logarithmic axis?
	Time             bool      // Time axis
	MinMode, MaxMode RangeMode // how to handel min and max
	TicSetting       TicSetting
	DataMin, DataMax float64 // actual values from data. if both zero: not calculated
	Min, Max         float64 // the min an d max of the xais
	TMin, TMax       *time.Time
	ShowLimits       bool
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
		r.ShowLimits = true
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
			if actNumTics > maxNumberOfTics { // TODO(vodo) this should never happen
				fmt.Printf("Switching to over next (%d > %d) delta from %s", actNumTics, maxNumberOfTics, td)
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

		// TODO(vodo) r.ShowLimits = true
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

type DataStyle struct {
	Symbol   int     // -1: no symbol, 0: auto, 1... fixed
	Line     int     // 0: no line, 1, solid, 2 dashed, 3 dotted, 4 dashdotted
	Size     float64 // 0: auto (1)
	Color    int     // index into palette
	Font     string
	FontSize int
	Alpha    float64
}


// Key placement
//          otl  otc  otr      
//         +--------------+ 
//     olt |itl  itc  itr | ort
//         |              |
//     olc |icl  icc  icr | ort
//         |              |
//     olb |ibl  ibc  ibr | orb
//         +--------------+ 
//        obl  obc  obr
type Key struct {
	Hide   bool
	Cols   int    // 
	Border int    // -1: off, 0: std, 1...:other styles
	Pos    string // "": itr
	// Width, Height int    // 0,0: auto
	X, Y    int
	Entries []KeyEntry
}

type KeyEntry struct {
	Symbol int
	Text   string
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


// Values and its std. implementation Real
type Value interface {
	XVal() float64
}

type Real float64

func (r Real) XVal() float64 { return float64(r) }


// XY-Values and its std. implementation Point
type XYValue interface {
	XVal() float64
	YVal() float64
}

type Point struct{ X, Y float64 }

func (p Point) XVal() float64            { return p.X }
func (p Point) YVal() float64            { return p.Y }
func (p Point) XErr() (float64, float64) { return math.NaN(), math.NaN() }
func (p Point) YErr() (float64, float64) { return math.NaN(), math.NaN() }


// XY-Value with error bars
type XYErrValue interface {
	XVal() float64
	YVal() float64
	XErr() (float64, float64)
	YErr() (float64, float64)
}
type EPoint struct {
	X, Y           float64
	DeltaX, DeltaY float64 // full range of x and y error, NaN for no errorbar
	OffX, OffY     float64 // offset of error range (must be < Delta)
}

func (p EPoint) XVal() float64 { return p.X }
func (p EPoint) YVal() float64 { return p.Y }
func (p EPoint) XErr() (float64, float64) {
	xl, _, xh, _ := p.BoundingBox()
	return xl, xh
}
func (p EPoint) YErr() (float64, float64) {
	_, yl, _, yh := p.BoundingBox()
	return yl, yh
}
func (p EPoint) BoundingBox() (xl, yl, xh, yh float64) { // bounding box
	xl, xh, yl, yh = p.X, p.X, p.Y, p.Y
	if !math.IsNaN(p.DeltaX) {
		xl -= p.DeltaX/2 - p.OffX
		xh += p.DeltaX/2 + p.OffX
	}
	if !math.IsNaN(p.DeltaY) {
		yl -= p.DeltaY/2 - p.OffY
		yh += p.DeltaY/2 + p.OffY
	}
	return
}

// Box in Boxplot
type Box struct {
	X, Avg, Med, Q1, Q3, Low, High float64
	Outliers                       []float64
}

func (p Box) XVal() float64 { return p.X }
func (p Box) YVal() float64 { return p.Med }
func (p Box) XErr() float64 { return p.Med - p.Q1 }
func (p Box) YErr() float64 { return p.Q3 - p.Med }


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

	width, leftm, height, topm = w-4, 2, h-1, 0
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

	numytics = h / 4

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
