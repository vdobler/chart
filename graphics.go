package chart

import (
	"fmt"
)


// Graphics is the interface all chart drivers have to implement
type Graphics interface {
	Dimensions() (int, int)        // character-width / height
	FontMetrics() (int, int)        // character-width / height

	Line(x0, y0, x1, y1 int, style DataStyle)                        // Draw line from (x0,y0) to (x1,y1)
	Symbol(x, y, s int, style DataStyle)                             // Put symnbol s at (x,y)
	Text(x, y int, t string, align string, rot int, style DataStyle) // align: [[tcb]][lcr]

	Begin() // start of chart drawing
	// All stuff is preprocessed: sanitized, clipped, strings formated, integer coords,
	// screen coordinates,
	XAxis(xr Range, ys, yms int)
	YAxis(yr Range, xs, xms int)
	Title(text string)

	/*
	Scatter(points []EPoint, style DataStyle) // Points, Lines and Line+Points
	Boxes(style DataStyle)                    // Boxplots
	Bars(style DataStyle)                     // any type of histogram
	Ring(style DataStyle)                     // 
	Key(entries []Key)
*/
	End() // Done, cleanup
}


// BasicGrapic is an interface of the most basic graphic primitives.
// Any type which implements BasicGraphics can use generic implementations
// of the Graphics methods.
type BasicGraphics interface {
	FontMetrics() (int, int)                                         // Return fontwidth and -height in pixel.
	Line(x0, y0, x1, y1 int, style DataStyle)                        // Draw line from (x0,y0) to (x1,y1)
	Symbol(x, y, s int, style DataStyle)                             // Put symnbol s at (x,y)
	Text(x, y int, t string, align string, rot int, style DataStyle) // align: [[tcb]][lcr]
	Style(element string) DataStyle                                  // retrieve style for element
}


// GenericAxis draws the axis r solely by graphic primitives of bg.
func GenericXAxis(bg BasicGraphics, rng Range, y, ym int) {
	_, fontheight := bg.FontMetrics()
	var ticLen int = 0
	if !rng.TicSetting.Hide {
		ticLen = min(10, max(4, (fontheight-1)/2))
	}

	// Axis itself, mirrord axis and zero
	xa, xe := rng.Data2Screen(rng.Min), rng.Data2Screen(rng.Max)
	bg.Line(xa, y, xe, y, bg.Style("axis"))
	if rng.TicSetting.Mirror >= 1 {
		bg.Line(xa, ym, xe, ym, bg.Style("maxis"))
	}
	if rng.ShowZero && rng.Min < 0 && rng.Max > 0 {
		z := rng.Data2Screen(0)
		bg.Line(z, y, z, ym, bg.Style("zero"))
	}

	// Axis label and range limits
	aly := y + 2*ticLen
	if !rng.TicSetting.Hide {
		aly += (3 * fontheight) / 2
	}
	if rng.ShowLimits {
		st := bg.Style("rangelimit")
		if rng.Time {
			bg.Text(xa, aly, rng.TMin.Format("2006-01-02 15:04:05"), "tl", 0, st)
			bg.Text(xe, aly, rng.TMax.Format("2006-01-02 15:04:05"), "tr", 0, st)
		} else {
			bg.Text(xa, aly, fmt.Sprintf("%g", rng.Min), "tl", 0, st)
			bg.Text(xe, aly, fmt.Sprintf("%g", rng.Max), "tr", 0, st)
		}
	}
	if rng.Label != "" { // draw label _after_ (=over) range limits
		bg.Text((xa+xe)/2, aly, "  "+rng.Label+"  ", "tc", 0, bg.Style("label"))
	}

	// Tics, tic labels an grid lines
	ticstyle := bg.Style("tic")
	for ticcnt, tic := range rng.Tics {
		x := rng.Data2Screen(tic.Pos)
		lx := rng.Data2Screen(tic.LabelPos)

		// Grid
		if ticcnt > 0 && ticcnt < len(rng.Tics)-1 && rng.TicSetting.Grid == 1 {
			fmt.Printf("Gridline at x=%d\n", x)
			bg.Line(x, y-1, x, ym+1, bg.Style("grid"))
		}

		// Tics
		fmt.Printf("y=%d  y-tl=%d  y+tl=%d\n", y, y-ticLen, y+ticLen)
		bg.Line(x, y-ticLen, x, y+ticLen, ticstyle)
		if rng.TicSetting.Mirror >= 2 {
			bg.Line(x, ym-ticLen, x, ym+ticLen, ticstyle)
		}
		if rng.Time && tic.Align == -1 {
			bg.Line(x, y+ticLen, x, y+2*ticLen, ticstyle)
			bg.Text(lx, y+2*ticLen, tic.Label, "tl", 0, ticstyle)
		} else {
			bg.Text(lx, y+2*ticLen, tic.Label, "tc", 0, ticstyle)
		}
	}

}

// GenericAxis draws the axis r solely by graphic primitives of bg.
func GenericYAxis(bg BasicGraphics, rng Range, x, xm int) {
	fontwidth, fontheight := bg.FontMetrics()
	var ticLen int = 0
	if !rng.TicSetting.Hide {
		ticLen = min(10, max(4, fontwidth))
	}

	// Axis itself, mirrord axis and zero
	ya, ye := rng.Data2Screen(rng.Min), rng.Data2Screen(rng.Max)
	bg.Line(x, ya, x, ye, bg.Style("axis"))
	if rng.TicSetting.Mirror >= 1 {
		bg.Line(xm, ya, xm, ye, bg.Style("maxis"))
	}
	if rng.ShowZero && rng.Min < 0 && rng.Max > 0 {
		z := rng.Data2Screen(0)
		bg.Line(x, z, xm, z, bg.Style("zero"))
	}

	// Label and axis ranges
	alx := 2 * fontheight
	if rng.ShowLimits {
		/* TODO
		st := bg.Style("rangelimit")
		if rng.Time {
			bg.Text(xa, aly, rng.TMin.Format("2006-01-02 15:04:05"), "tl", 0, st)
			bg.Text(xe, aly, rng.TMax.Format("2006-01-02 15:04:05"), "tr", 0, st)
		} else {
			bg.Text(xa, aly, fmt.Sprintf("%g", rng.Min), "tl", 0, st)
			bg.Text(xe, aly, fmt.Sprintf("%g", rng.Max), "tr", 0, st)
		}
		 */
	}
	if rng.Label != "" {
		y := (ya+ye)/2 
		bg.Text(alx, y, rng.Label, "bc", 90, bg.Style("label"))
	}

	// Tics, tic labels and grid lines
	ticstyle := bg.Style("tic")
	for _, tic := range rng.Tics {
		y := rng.Data2Screen(tic.Pos)
		ly := rng.Data2Screen(tic.LabelPos)
		bg.Line(x-ticLen, y, x+ticLen, y, ticstyle)
		if rng.TicSetting.Mirror >= 2 {
			bg.Line(xm-ticLen, y, xm+ticLen, y, ticstyle)
		}
		if rng.Time && tic.Align == 0 { // centered tic
			bg.Line(x-2*ticLen, y, x+ticLen, y, ticstyle)
			bg.Text(x-ticLen, ly, tic.Label, "cr", 90, ticstyle)
		} else {
			bg.Text(x-2*ticLen, ly, tic.Label, "cr", 0, ticstyle)
		}
	}
}