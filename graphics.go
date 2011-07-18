package chart

import (
	"fmt"
)


type Graphics interface {
	AxisOrientation() (bool, bool) // yield axis orientation. true: normal
	FontMetric() (int, int)        // character-width / height

	Begin() // start of chart drawing
	// All stuff is preprocessed: sanitized, clipped, strings formated, integer coords,
	// screen coordinates,
	XAxis(xr Range, ys, yms int)
	YAxis(r Range)
	Title(text string)
	Scatter(points []EPoint, style DataStyle) // Points, Lines and Line+Points
	Boxes(boxes []Box, style DataStyle)       // Boxplots
	Bars(bars []Bar, style DataStyle)         // any type of histogram
	Ring(segments []Segment, style DataStyle) // 
	Key(entries []Key)

	End() // Done, cleanup
}


// BasicGrapic is an interface of the most basic graphic primitives
type BasicGraphics interface {
	Line(x0, y0, x1, y1, int, style DataStyle)
	Symbol(x, y, s int, style DataStyle)
	Text(x, y int, t string, align string, rot int, style DataStyle) // align: [[tcb]][lcr]
	Style(element string) DataStyle                                  // retrieve style for element
}


// GenericAxis draws the axis r solely by graphic primitives of bg.

func GenericXAxis(bg *BasicGraphics, r Range, y, ym int) {
	var ticLen int = 0
	if !xrange.TicSetting.Hide {
		ticLen = 5
	}

	xa, xe := xrange.Data2Screen(xrange.Min), xrange.Data2Screen(xrange.Max)
	bg.Line(xa, y, xe, y, bg.Style("axis"))
	if mirror >= 1 {
		bg.Line(xa, ym, xe, ym, bg.Style("maxis"))
	}
	if xrange.ShowZero && xrange.Min < 0 && xrange.Max > 0 {
		z := xrange.Data2Screen(0)
		bg.Line(z, y, z, y1, bg.Style("zero"))
	}

	if label != "" {
		yy := y + 5 + ticlen
		if !xrange.TicSetting.Hide {
			yy += 3*fontheight/2 + ticLen
		}
		bg.Text((xa+xe)/2, yy, label, "tc")
	}

	for ticcnt, tic := range xrange.Tics {
		x := xrange.Data2Screen(tic.Pos)
		lx := xrange.Data2Screen(tic.LabelPos)
		if ticcnt > 0 && ticcnt < len(xrange.Tics)-1 && xrange.TicSetting.Grid == 1 {
			bg.Line(x, y-1, x, ym+1, bg.Style("tic"))
		}
		bg.Line(x, y-ticLen, x, y+ticLen, bg.Style("tic"))
		if mirror >= 2 {
			bg.Line(x, ym-ticLen, x, ym+ticLen, bg.Style("tic"))
		}
		if xrange.Time {
			chart.Line(x, y+ticLen, x, y+2*ticLen, "stroke:black; stroke-width:2")
			if tic.Align == -1 {
				chart.Text(lx, y+fontheight+ticLen, tic.Label, "text-anchor:left")
			} else {
				chart.Text(lx, y+fontheight+ticLen, tic.Label, "text-anchor:middle")
			}
		} else {
			chart.Text(lx, y+fontheight+ticLen, tic.Label, "text-anchor:middle")
		}
	}
	if xrange.ShowLimits {
		/*
		 if xrange.Time {
		 tb.Text(xa, y+2, xrange.TMin.Format("2006-01-02 15:04:05"), -1)
		 tb.Text(xe, y+2, xrange.TMax.Format("2006-01-02 15:04:05"), 1)
		 } else {
		 tb.Text(xa, y+2, fmt.Sprintf("%g", xrange.Min), -1)
		 tb.Text(xe, y+2, fmt.Sprintf("%g", xrange.Max), 1)
		 }
		*/
	}

}
