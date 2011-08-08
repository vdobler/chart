package chart

import (
	"fmt"
	"math"
)

// BasicGrapic is an interface of the most basic graphic primitives.
// Any type which implements BasicGraphics can use generic implementations
// of the Graphics methods.
type BasicGraphics interface {
	Style(element string) DataStyle                         // retrieve style for element
	Font(element string) Font                               // retrieve font for element
	FontMetrics(font Font) (fw float32, fh int, mono bool)  // Return fontwidth and -height in pixel
	TextLen(t string, font Font) int                        // length=width of t in screen units
	Line(x0, y0, x1, y1 int, style DataStyle)               // Draw line from (x0,y0) to (x1,y1)
	Symbol(x, y, s int, style DataStyle)                    // Put symnbol s at (x,y)
	Text(x, y int, t string, align string, rot int, f Font) // align: [[tcb]][lcr]
	Rect(x, y, w, h int, style DataStyle)                   // draw (w x h) rectangle at (x,y)
	Wedge(x, y, r int, phi, psi float64, style DataStyle)   // draw pie from phi to psi centered at (x,y) with radius r
}


// Graphics is the interface all chart drivers have to implement
type Graphics interface {
	BasicGraphics

	Dimensions() (int, int) // character-width / height

	Begin() // start of chart drawing
	// All stuff is preprocessed: sanitized, clipped, strings formated, integer coords,
	// screen coordinates,
	XAxis(xr Range, ys, yms int) // Draw x axis xr at screen position ys (and yms if mirrored)
	YAxis(yr Range, xs, xms int) // Same for y axis.
	Title(text string)           // Draw title onto chart

	Scatter(points []EPoint, plotstyle PlotStyle, style DataStyle) // Points, Lines and Line+Points
	Boxes(boxes []Box, width int, style DataStyle)                 // Boxplots
	Bars(bars []Barinfo, style DataStyle)                          // any type of histogram
	/*
		Ring(style DataStyle)                     // 
	*/
	Key(x, y int, key Key) // place key at x,y
	End()                  // Done, cleanup
}

type Barinfo struct {
	x, y   int    // (x,y) of top left corner; 
	w, h   int    // width and heigt
	t, tp  string // label text and text position '[oi][tblr]' or 'c'
	f      Font   // font of text
}


func GenericTextLen(bg BasicGraphics, t string, font Font) (width int) {
	// TODO: how handle newlines?  same way like Text does
	fw, _, mono := bg.FontMetrics(font)
	if mono {
		for _ = range t {
			width++
		}
		width = int(float32(width)*fw + 0.5)
	} else {
		var length float32
		for _, rune := range t {
			if w, ok := CharacterWidth[rune]; ok {
				length += w
			} else {
				length += 20 // save above average
			}
		}
		length /= averageCharacterWidth
		length *= fw
		width = int(length + 0.5)
	}
	return
}

// GenericRect draws a rectangle of size w x h at (x,y).  Drawing is done
// by simple lines only.
func GenericRect(bg BasicGraphics, x, y, w, h int, style DataStyle) {
	if style.FillColor != "" {
		fs := DataStyle{LineWidth: 1, LineColor: style.FillColor, LineStyle: SolidLine, Alpha: 0}
		for i := 1; i < h-1; i++ {
			bg.Line(x+1, y+i, x+w-1, y+i, fs)
		}
	}
	bg.Line(x, y, x+w, y, style)
	bg.Line(x+w, y, x+w, y+h, style)
	bg.Line(x+w, y+h, x, y+h, style)
	bg.Line(x, y+h, x, y, style)
}

// GenericAxis draws the axis r solely by graphic primitives of bg.
func GenericXAxis(bg BasicGraphics, rng Range, y, ym int) {
	_, fontheight, _ := bg.FontMetrics(bg.Font("label"))
	ticfont := bg.Font("tic")
	var ticLen int = 0
	if !rng.TicSetting.Hide {
		ticLen = min(10, max(4, fontheight/2))
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
		f := bg.Font("rangelimit")
		if rng.Time {
			bg.Text(xa, aly, rng.TMin.Format("2006-01-02 15:04:05"), "tl", 0, f)
			bg.Text(xe, aly, rng.TMax.Format("2006-01-02 15:04:05"), "tr", 0, f)
		} else {
			bg.Text(xa, aly, fmt.Sprintf("%g", rng.Min), "tl", 0, f)
			bg.Text(xe, aly, fmt.Sprintf("%g", rng.Max), "tr", 0, f)
		}
	}
	if rng.Label != "" { // draw label _after_ (=over) range limits
		bg.Text((xa+xe)/2, aly, "  "+rng.Label+"  ", "tc", 0, bg.Font("label"))
	}

	if !rng.TicSetting.Hide {
		// Tics, tic labels an grid lines
		ticstyle := bg.Style("tic")
		for ticcnt, tic := range rng.Tics {
			x := rng.Data2Screen(tic.Pos)
			lx := rng.Data2Screen(tic.LabelPos)

			// Grid
			if ticcnt > 0 && ticcnt < len(rng.Tics)-1 && rng.TicSetting.Grid == 1 {
				// fmt.Printf("Gridline at x=%d\n", x)
				bg.Line(x, y-1, x, ym+1, bg.Style("grid"))
			}

			// Tics
			// fmt.Printf("y=%d  y-tl=%d  y+tl=%d\n", y, y-ticLen, y+ticLen)
			bg.Line(x, y-ticLen, x, y+ticLen, ticstyle)
			if rng.TicSetting.Mirror >= 2 {
				bg.Line(x, ym-ticLen, x, ym+ticLen, ticstyle)
			}
			if rng.Time && tic.Align == -1 {
				bg.Line(x, y+ticLen, x, y+2*ticLen, ticstyle)
				bg.Text(lx, y+2*ticLen, tic.Label, "tl", 0, ticfont)
			} else {
				bg.Text(lx, y+2*ticLen, tic.Label, "tc", 0, ticfont)
			}
		}
	}
}

// GenericAxis draws the axis r solely by graphic primitives of bg.
func GenericYAxis(bg BasicGraphics, rng Range, x, xm int) {
	_, fontheight, _ := bg.FontMetrics(bg.Font("label"))
	var ticLen int = 0
	if !rng.TicSetting.Hide {
		ticLen = min(10, max(4, fontheight/2))
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
		y := (ya + ye) / 2
		bg.Text(alx, y, rng.Label, "bc", 90, bg.Font("label"))
	}

	if !rng.TicSetting.Hide {
		// Tics, tic labels and grid lines
		ticstyle := bg.Style("tic")
		ticfont := bg.Font("tic")
		for ticcnt, tic := range rng.Tics {
			y := rng.Data2Screen(tic.Pos)
			ly := rng.Data2Screen(tic.LabelPos)

			// Grid
			if ticcnt > 0 && ticcnt < len(rng.Tics)-1 && rng.TicSetting.Grid == 1 {
				// fmt.Printf("Gridline at x=%d\n", x)
				bg.Line(x+1, y, xm-1, y, bg.Style("grid"))
			}

			// Tics
			bg.Line(x-ticLen, y, x+ticLen, y, ticstyle)
			if rng.TicSetting.Mirror >= 2 {
				bg.Line(xm-ticLen, y, xm+ticLen, y, ticstyle)
			}
			if rng.Time && tic.Align == 0 { // centered tic
				bg.Line(x-2*ticLen, y, x+ticLen, y, ticstyle)
				bg.Text(x-ticLen, ly, tic.Label, "cr", 90, ticfont)
			} else {
				bg.Text(x-2*ticLen, ly, tic.Label, "cr", 0, ticfont)
			}
		}
	}
}


// GenericScatter draws the given points according to style.
// style.FillColor is used as color of error bars and style.FontSize is used
// as the length of the endmarks of the error bars. Both have suitable defaults
// if the FontXyz are not set. Point coordinates and errors must be provided 
// in screen coordinates.
func GenericScatter(bg BasicGraphics, points []EPoint, plotstyle PlotStyle, style DataStyle) {
	// First pass: Error bars
	for _, p := range points {
		ebs := style
		ebs.LineColor, ebs.LineWidth = ebs.FillColor, 1
		if ebs.LineColor == "" {
			ebs.LineColor = "#404040"
		}
		if ebs.LineWidth == 0 {
			ebs.LineWidth = 1
		}
		xl, yl, xh, yh := p.boundingBox()
		// fmt.Printf("Draw %d: %f %f-%f\n", i, p.DeltaX, xl,xh)
		if !math.IsNaN(p.DeltaX) {
			bg.Line(int(xl), int(p.Y), int(xh), int(p.Y), ebs)
		}
		if !math.IsNaN(p.DeltaY) {
			bg.Line(int(p.X), int(yl), int(p.X), int(yh), ebs)
		}
	}

	// Second pass: Line
	if (plotstyle&PlotStyleLines) != 0 && len(points) > 0 {
		lastx, lasty := points[0].X, points[0].Y
		for i := 1; i < len(points); i++ {
			x, y := points[i].X, points[i].Y
			bg.Line(int(lastx), int(lasty), int(x), int(y), style)
			lastx, lasty = x, y
		}
	}

	// Third pass: symbols
	if (plotstyle&PlotStylePoints) != 0 && len(points) != 0 {
		for _, p := range points {
			bg.Symbol(int(p.X), int(p.Y), style.Symbol, style)
		}
	}
}

// GenericBoxes draws box plots. (Default implementation for box plots).
func GenericBoxes(bg BasicGraphics, boxes []Box, width int, style DataStyle) {
	if width%2 == 0 {
		width += 1
	}
	hbw := (width - 1) / 2
	for _, d := range boxes {
		x := int(d.X)
		q1, q3 := int(d.Q1), int(d.Q3)

		bg.Rect(x-hbw, q1, width, q3-q1, style)
		if !math.IsNaN(d.Med) {
			med := int(d.Med)
			bg.Line(x-hbw, med, x+hbw, med, style)
		}

		if !math.IsNaN(d.Avg) {
			bg.Symbol(x, int(d.Avg), style.Symbol, style)
		}

		if !math.IsNaN(d.High) {
			bg.Line(x, q3, x, int(d.High), style)
		}

		if !math.IsNaN(d.Low) {
			bg.Line(x, q1, x, int(d.Low), style)
		}

		for _, y := range d.Outliers {
			bg.Symbol(x, int(y), style.Symbol, style)
		}

	}

}

// TODO: Is Bars and Generic Bars useful at all? Replaceable by rect?
func GenericBars(bg BasicGraphics, bars []Barinfo, style DataStyle) {
	for _, b := range bars {
		bg.Rect(b.x, b.y, b.w, b.h, style)
		if b.t != "" {
			var tx, ty int
			var a string
			_, fh, _ := bg.FontMetrics(b.f)
			if fh > 1 {
				fh /= 2
			}
			switch b.tp {
			case "ot": 
				tx, ty, a = b.x + b.w/2, b.y - fh, "bc"
			case "it": 
				tx, ty, a = b.x + b.w/2, b.y + fh, "tc"
			case "ib": 
				tx, ty, a = b.x + b.w/2, b.y + b.h - fh, "bc"
			case "ob": 
				tx, ty, a = b.x + b.w/2, b.y + b.h + fh, "tc"
			case "ol": 
				tx, ty, a = b.x - fh, b.y + b.h/2, "cr"
			case "il": 
				tx, ty, a = b.x + fh, b.y + b.h/2, "cl"
			case "or": 
				tx, ty, a = b.x + b.w + fh, b.y + b.h/2, "cl"
			case "ir": 
				tx, ty, a = b.x + b.w - fh, b.y + b.h/2, "cr"
			default: 
				tx, ty, a = b.x + b.w/2, b.y + b.h/2, "cc"
				
			}
				
				
			bg.Text(tx, ty, b.t, a, 0, b.f)
		}
	}
}

// GenericWedge draws a pie/wedge just by lines
func GenericWedge(bg BasicGraphics, x, y, r int, phi, psi float64, style DataStyle) {

	xa, ya := int(math.Cos(phi)*float64(r))+x, int(math.Sin(phi)*float64(r))+y
	xc, yc := int(math.Cos(psi)*float64(r))+x, int(math.Sin(psi)*float64(r))+y
	rf := float64(r)

	if math.Fabs(phi-psi) >= 4*math.Pi {
		phi, psi = 0, 2*math.Pi
	} else {
		bg.Line(x, y, xa, ya, style)
		bg.Line(x, y, xc, yc, style)
	}

	if style.FillColor != "" {
		delta := 1 / (4 * rf)
		ls := DataStyle{LineColor: style.FillColor, LineWidth: 2, Symbol: style.Symbol}
		for a := phi; a <= psi; a += delta {
			xr, yr := int(math.Cos(a)*rf)+x, int(math.Sin(a)*rf)+y
			bg.Line(x, y, xr, yr, ls)
		}
	}

	var xb, yb int
	for ; phi < psi; phi += 0.1 { // aproximate circle by 62-corner
		xb, yb = int(math.Cos(phi)*rf)+x, int(math.Sin(phi)*rf)+y
		bg.Line(xa, ya, xb, yb, style)
		xa, ya = xb, yb
	}
	bg.Line(xb, yb, xc, yc, style)
}
