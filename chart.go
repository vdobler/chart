package chart

import (
	"fmt"
	"math"
	//	"os"
	//	"strings"
)


type RangeExpansion int

const (
	ExpandNextTic = iota
	ExpandToTic
	ExpandTight
	ExpandABit
)

type RangeMode struct {
	// false: autoscaling; true: use Value as fixed setting
	Fixed bool
	// false: unconstrained autoscaling; true: use Lower and Upper as limits
	Constrained bool
	// one if ExpandNextTic, ExpandTight, ExpandABit
	Expand int
	// see above
	Value, Lower, Upper float64
}


type Tics struct {
	Hide               bool    // dont show
	First, Last, Delta float64 // Delta is factor for log axis
	Minor              int     // 0: off, 1: clever, >1: number of intervalls
}


type Range struct {
	// logarithmic axis?
	Log bool
	// date axis
	Date bool
	// how to handel min and max
	MinMode, MaxMode RangeMode
	// actual values from data. if both zero: not calculated
	DataMin, DataMax float64
	// the min an d max of the xais
	Min, Max float64
	Tics     Tics
	// Convert data coordnate to screen
	Data2Screen func(float64) int
	// Convert screen coordinate to data
	Screen2Data func(int) float64
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

func almostEqual(a, b float64) bool {
	rd := math.Fabs((a - b) / (a + b))
	return rd < 1e-5
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

// ticdistances: 1 2 or 5 * 10^n to have min (2) few (3, 4), 
// some (5,6), several (7,8,9) or many (10) 
func (r *Range) Setup(numberOfTics, sWidth, sOffset int, revert bool) {
	if numberOfTics <= 1 { numberOfTics = 2 }
	delta := (r.DataMax - r.DataMin) / float64(numberOfTics-1)
	if delta == 0 { delta = 1 }
	de := math.Floor(math.Log10(delta))
	// fmt.Printf(":: %f, %f, %d \n", delta, de, int(de))
	f := delta / math.Pow10(int(de))
	switch {
	case f < 2:
		f = 1
	case f < 4:
		f = 2
	case f < 9:
		f = 5
	default:
		f = 1
		de += 1
	}
	delta = f * math.Pow10(int(de))
	// fmt.Printf("delta = %f,  f = %f\n", delta, f)

	r.Min = Bound(r.MinMode, r.DataMin, delta, false)
	r.Max = Bound(r.MaxMode, r.DataMax, delta, true)
	r.Tics.First = delta * math.Ceil(r.Min/delta /* - 0.001*delta */ )
	r.Tics.Last = delta * math.Floor(r.Max/delta /* + 0.001*delta */ )
	r.Tics.Delta = delta
	//	if first == last { last += delta } // TODO
	fmt.Printf("Range: (%g,%g) --> (%g,%g), Tic-Delta %g from %g to %g\n", r.DataMin, r.DataMax, r.Min, r.Max, delta, r.Tics.First, r.Tics.Last)

	if !revert {
		r.Data2Screen = func(x float64) int {
			return int(math.Floor(float64(sWidth)*(x-r.Min)/(r.Max-r.Min))) + sOffset
		}
		r.Screen2Data = func(x int) float64 {
			return (r.Max-r.Min)*float64(x-sOffset)/float64(sWidth) + r.Min
		}
	} else {
		r.Data2Screen = func(x float64) int {
			return sWidth-int(math.Floor(float64(sWidth)*(x-r.Min)/(r.Max-r.Min))) + sOffset
		}
		r.Screen2Data = func(x int) float64 {
			return (r.Max-r.Min)*float64(-x+sOffset+sWidth)/float64(sWidth) + r.Min
		}

	}

}

type KeyEntry struct {
	Symbol int
	Text string
}

type DataStyle struct {
	Symbol int   // -1: no symbol, 0: auto, 1... fixed
	Line int     // 0: no line, 1, solid, 2 dashed, 3 dotted
	Size float64 // 0: auto (1)
	Color int    // index into palette
}

type Point struct {
	X, Y float64
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
	Border        int    // -1: off, 0: std, 1...:other styles
	Pos           string // "": itr
	// Width, Height int    // 0,0: auto
	X, Y int
}


func min(a, b int) int { 
	if a<b { 
		return a 
	}
	return b
}

func max(a, b int) int { 
	if a>b { 
		return a 
	}
	return b
}

func abs(a int) int { 
	if a < 0 {
		return -a
	}
	return a
}

func clip(x, l, u int) int {
	if x < min(l,u) { return l }
	if x > max(l,u) { return u }
	return x
}
