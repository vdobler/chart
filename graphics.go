package chart

import (
	"fmt"
	"math"
)

// BasicGrapic is an interface of the most basic graphic primitives.
// Any type which implements BasicGraphics can use generic implementations
// of the Graphics methods.
type BasicGraphics interface {
	Style(element string) Style                             // Retrieve style for element
	Font(element string) Font                               // Retrieve font for element
	FontMetrics(font Font) (fw float32, fh int, mono bool)  // Return fontwidth and -height in pixel
	TextLen(t string, font Font) int                        // Length=width of t in screen units if set on font 
	Line(x0, y0, x1, y1 int, style Style)                   // Draw line from (x0,y0) to (x1,y1)
	Symbol(x, y, s int, style Style)                        // Put symbol s at (x,y)
	Text(x, y int, t string, align string, rot int, f Font) // Put t at (x,y) rotated by rot aligned [[tcb]][lcr]
	Rect(x, y, w, h int, style Style)                       // Draw (w x h) rectangle at (x,y)
}


// Graphics is the interface all chart drivers have to implement
type Graphics interface {
	BasicGraphics

	Dimensions() (int, int) // character-width / height

	Begin() // start of chart drawing
	End()   // Done, cleanup

	// All stuff is preprocessed: sanitized, clipped, strings formated, integer coords,
	// screen coordinates,
	XAxis(xr Range, ys, yms int) // Draw x axis xr at screen position ys (and yms if mirrored)
	YAxis(yr Range, xs, xms int) // Same for y axis.
	Title(text string)           // Draw title onto chart, box if l,r,y != 0

	Scatter(points []EPoint, plotstyle PlotStyle, style Style) // Points, Lines and Line+Points
	Boxes(boxes []Box, width int, style Style)                 // Boxplots
	Bars(bars []Barinfo, style Style)                          // any type of histogram/bars
	Rings(wedeges []Wedgeinfo, x, y, ro, ri int)               // Pie/ring diagram elements

	Key(x, y int, key Key) // place key at x,y
}

type Barinfo struct {
	x, y  int    // (x,y) of top left corner; 
	w, h  int    // width and heigt
	t, tp string // label text and text position '[oi][tblr]' or 'c'
	f     Font   // font of text
}

type Wedgeinfo struct {
	Phi, Psi float64 // Start and ende of wedge. Fuill circle if |phi-psi| > 4pi
	Text, Tp string  // label text and text position: [ico]
	Style    Style   // style of this wedge
	Font     Font    // font of text
	Shift    int     // Highlighting of wedge 
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
func GenericRect(bg BasicGraphics, x, y, w, h int, style Style) {
	if style.FillColor != "" {
		fs := Style{LineWidth: 1, LineColor: style.FillColor, LineStyle: SolidLine, Alpha: 0}
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
		ticLen = imin(10, imax(4, fontheight/2))
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
		ticLen = imin(10, imax(4, fontheight/2))
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
func GenericScatter(bg BasicGraphics, points []EPoint, plotstyle PlotStyle, style Style) {

	// First pass: Error bars
	ebs := style
	ebs.LineColor, ebs.LineWidth, ebs.LineStyle = ebs.FillColor, 1, SolidLine
	if ebs.LineColor == "" {
		ebs.LineColor = "#404040"
	}
	if ebs.LineWidth == 0 {
		ebs.LineWidth = 1
	}
	for _, p := range points {
		xl, yl, xh, yh := p.BoundingBox()
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
		lastx, lasty := int(points[0].X), int(points[0].Y)
		for i := 1; i < len(points); i++ {
			x, y := int(points[i].X), int(points[i].Y)
			bg.Line(lastx, lasty, x, y, style)
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
func GenericBoxes(bg BasicGraphics, boxes []Box, width int, style Style) {
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
func GenericBars(bg BasicGraphics, bars []Barinfo, style Style) {
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
				tx, ty, a = b.x+b.w/2, b.y-fh, "bc"
			case "it":
				tx, ty, a = b.x+b.w/2, b.y+fh, "tc"
			case "ib":
				tx, ty, a = b.x+b.w/2, b.y+b.h-fh, "bc"
			case "ob":
				tx, ty, a = b.x+b.w/2, b.y+b.h+fh, "tc"
			case "ol":
				tx, ty, a = b.x-fh, b.y+b.h/2, "cr"
			case "il":
				tx, ty, a = b.x+fh, b.y+b.h/2, "cl"
			case "or":
				tx, ty, a = b.x+b.w+fh, b.y+b.h/2, "cl"
			case "ir":
				tx, ty, a = b.x+b.w-fh, b.y+b.h/2, "cr"
			default:
				tx, ty, a = b.x+b.w/2, b.y+b.h/2, "cc"

			}

			bg.Text(tx, ty, b.t, a, 0, b.f)
		}
	}
}

// GenericWedge draws a pie/wedge just by lines
func GenericWedge(bg BasicGraphics, x, y, ro, ri int, phi, psi, ecc float64, style Style) {
	for phi < 0 {
		phi += 2 * math.Pi
	}
	for psi < 0 {
		psi += 2 * math.Pi
	}
	for phi >= 2*math.Pi {
		phi -= 2 * math.Pi
	}
	for psi >= 2*math.Pi {
		psi -= 2 * math.Pi
	}
	debug.Printf("GenericWedge centered at (%d,%d) from %.1f° to %.1f°, radius %d/%d (e=%.2f)",
		x, y, 180*phi/math.Pi, 180*psi/math.Pi, ro, ri, ecc)

	if ri > ro {
		panic("ri > ro is not possible")
	}

	if style.FillColor != "" {
		fillWedge(bg, x, y, ro, ri, phi, psi, ecc, style)
	}

	roe, rof := float64(ro)*ecc, float64(ro)
	rie, rif := float64(ri)*ecc, float64(ri)
	xa, ya := int(math.Cos(phi)*roe)+x, y-int(math.Sin(phi)*rof)
	xc, yc := int(math.Cos(psi)*roe)+x, y-int(math.Sin(psi)*rof)
	xai, yai := int(math.Cos(phi)*rie)+x, y-int(math.Sin(phi)*rif)
	xci, yci := int(math.Cos(psi)*rie)+x, y-int(math.Sin(psi)*rif)

	if math.Fabs(phi-psi) >= 4*math.Pi {
		phi, psi = 0, 2*math.Pi
	} else {
		if ri > 0 {
			bg.Line(xai, yai, xa, ya, style)
			bg.Line(xci, yci, xc, yc, style)
		} else {
			bg.Line(x, y, xa, ya, style)
			bg.Line(x, y, xc, yc, style)
		}
	}

	var xb, yb int
	exit := phi < psi
	for rho := phi; !exit || rho < psi; rho += 0.05 { // aproximate circle by more than 120 corners polygon
		if rho >= 2*math.Pi {
			exit = true
			rho -= 2 * math.Pi
		}
		xb, yb = int(math.Cos(rho)*roe)+x, y-int(math.Sin(rho)*rof)
		bg.Line(xa, ya, xb, yb, style)
		xa, ya = xb, yb
	}
	bg.Line(xb, yb, xc, yc, style)

	if ri > 0 {
		exit := phi < psi
		for rho := phi; !exit || rho < psi; rho += 0.1 { // aproximate circle by more than 60 corner polygon
			if rho >= 2*math.Pi {
				exit = true
				rho -= 2 * math.Pi
			}
			xb, yb = int(math.Cos(rho)*rie)+x, y-int(math.Sin(rho)*rif)
			bg.Line(xai, yai, xb, yb, style)
			xai, yai = xb, yb
		}
		bg.Line(xb, yb, xci, yci, style)

	}
}

// Fill wedge with center (xi,yi), radius ri from alpha to beta with style.
// Precondition:  0 <= beta < alpha < pi/2
func fillQuarterWedge(bg BasicGraphics, xi, yi, ri int, alpha, beta, e float64, style Style, quadrant int) {
	if alpha < beta {
		// debug.Printf("Swaping alpha and beta")
		alpha, beta = beta, alpha
	}
	// debug.Printf("fillQuaterWedge from %.1f to %.1f radius %d in quadrant %d.",	180*alpha/math.Pi, 180*beta/math.Pi, ri, quadrant)
	r := float64(ri)

	ta, tb := math.Tan(alpha), math.Tan(beta)
	for y := int(r * math.Sin(alpha)); y >= 0; y-- {
		yf := float64(y)
		x0 := yf / ta
		x1 := yf / tb
		x2 := math.Sqrt(r*r - yf*yf)
		// debug.Printf("y=%d  x0=%.2f    x1=%.2f    x2=%.2f  border=%t", y, x0, x1, x2, (x2<x1))  
		if math.IsNaN(x1) || x2 < x1 {
			x1 = x2
		}

		var xx0, xx1, yy int
		switch quadrant {
		case 0:
			xx0 = int(x0*e+0.5) + xi
			xx1 = int(x1*e-0.5) + xi
			yy = yi - y
		case 3:
			xx0 = int(x0*e+0.5) + xi
			xx1 = int(x1*e-0.5) + xi
			yy = yi + y
		case 2:
			xx0 = xi - int(x0*e+0.5)
			xx1 = xi - int(x1*e-0.5)
			yy = yi + y
		case 1:
			xx0 = xi - int(x0*e+0.5)
			xx1 = xi - int(x1*e-0.5)
			yy = yi - y
		default:
			panic("No such quadrant.")
		}
		// debug.Printf("Line %d,%d to %d,%d", xx0,yy, xx1,yy)
		bg.Line(xx0, yy, xx1, yy, style)
	}
}

func quadrant(w float64) int {
	return int(math.Floor(2 * w / math.Pi))
}

func mapQ(w float64, q int) float64 {
	switch q {
	case 0:
		return w
	case 1:
		return math.Pi - w
	case 2:
		return w - math.Pi
	case 3:
		return 2*math.Pi - w
	default:
		panic("No such quadrant")
	}
	return w
}

// Fill wedge with center (xi,yi), radius ri from alpha to beta with style.
// Any combination of phi, psi allowed as long 0 <= phi < psi < 2pi.
func fillWedge(bg BasicGraphics, xi, yi, ro, ri int, phi, psi, epsilon float64, style Style) {
	// ls := Style{LineColor: style.FillColor, LineWidth: 1, Symbol: style.Symbol}

	qPhi := quadrant(phi)
	qPsi := quadrant(psi)
	// debug.Printf("fillWedge from %.1f (%d) to %.1f (%d).", 180*phi/math.Pi, qPhi, 180*psi/math.Pi, qPsi)

	// prepare styles for filling
	style.LineColor = style.FillColor
	style.LineWidth = 1
	style.LineStyle = SolidLine
	blank := Style{Symbol: ' ', LineColor: "#ffffff", Alpha: 1}

	for qPhi != qPsi {
		// debug.Printf("qPhi = %d", qPhi)
		w := float64(qPhi+1) * math.Pi / 2
		if math.Fabs(w-phi) > 0.01 {
			fillQuarterWedge(bg, xi, yi, ro, mapQ(phi, qPhi), mapQ(w, qPhi), epsilon, style, qPhi)
			if ri > 0 {
				fillQuarterWedge(bg, xi, yi, ri, mapQ(phi, qPhi), mapQ(w, qPhi), epsilon, blank, qPhi)
			}
		}
		phi = w
		qPhi++
		if qPhi == 4 {
			// debug.Printf("Wrapped phi around")
			phi, qPhi = 0, 0
		}
	}
	if phi != psi {
		// debug.Printf("Last wedge")
		fillQuarterWedge(bg, xi, yi, ro, mapQ(phi, qPhi), mapQ(psi, qPhi), epsilon, style, qPhi)
		if ri > 0 {
			fillQuarterWedge(bg, xi, yi, ri, mapQ(phi, qPhi), mapQ(psi, qPhi), epsilon, blank, qPhi)
		}
	}
}


func GenericRings(bg BasicGraphics, wedges []Wedgeinfo, x, y, ro, ri int, eccentricity float64) {
	debug.Printf("GenericRings with %d wedges center %d,%d, radii %d/%d,  ecc=%.3f)",
		len(wedges), x, y, ro, ri, eccentricity)

	for _, w := range wedges {

		// Correct center
		p := 0.4 * float64(w.Style.LineWidth+w.Shift)

		// cphi, sphi := math.Cos(w.Phi), math.Sin(w.Phi)
		// cpsi, spsi := math.Cos(w.Psi), math.Sin(w.Psi)
		a := math.Sin((w.Psi - w.Phi) / 2)
		dx, dy := p*math.Cos((w.Phi+w.Psi)/2)/a, p*math.Sin((w.Phi+w.Psi)/2)/a
		debug.Printf("Center adjustment (lw=%d, p=%.2f), for wedge %d°-%d° of (%.1f,%.1f)",
			w.Style.LineWidth, p, int(180*w.Phi/math.Pi), int(180*w.Psi/math.Pi), dx, dy)
		xi, yi := x+int(dx+0.5), y-int(dy+0.5)
		GenericWedge(bg, xi, yi, ro, ri, w.Phi, w.Psi, eccentricity, w.Style)

		if w.Text != "" {
			_, fh, _ := bg.FontMetrics(w.Font)
			fh += 0
			alpha := (w.Phi + w.Psi) / 2
			var rt int
			if ri > 0 {
				rt = (ri + ro) / 2
			} else {
				rt = ro - 3*fh
				if rt <= ro/2 {
					rt = ro - 2*fh
				}
			}
			// debug.Printf("Text %s at %d° r=%d", w.Text, int(180*alpha/math.Pi), rt)
			tx := int(float64(rt)*math.Cos(alpha)*eccentricity+0.5) + x
			ty := y - int(float64(rt)*math.Sin(alpha)+0.5)

			bg.Text(tx, ty, w.Text, "cc", 0, w.Font)
		}

	}

	/***************
	var d string
	p := 0.4 * float64(w.Style.LineWidth)
	cphi, sphi := math.Cos(w.Phi), math.Sin(w.Phi)
	cpsi, spsi := math.Cos(w.Psi), math.Sin(w.Psi)

	if ri <= 0 {
		// real wedge drawn as center -> outer radius -> arc -> closed to center
		rf := float64(ro)
		a := math.Sin((w.Psi - w.Phi) / 2)
		dx, dy := p*math.Cos((w.Phi+w.Psi)/2)/a, p*math.Sin((w.Phi+w.Psi)/2)/a
		d = fmt.Sprintf("M %d,%d ", x+int(dx+0.5), y+int(dy+0.5))

		dx, dy = p*math.Cos(w.Phi+math.Pi/2), p*math.Sin(w.Phi+math.Pi/2)
		d += fmt.Sprintf("L %d,%d ", int(rf*cphi+0.5+dx)+x, int(rf*sphi+0.5+dy)+y)

		dx, dy = p*math.Cos(w.Psi-math.Pi/2), p*math.Sin(w.Psi-math.Pi/2)
		d += fmt.Sprintf("A %d,%d 0 0 1 %d,%d ", ro, ro, int(rf*cpsi+0.5+dx)+x, int(rf*spsi+0.5+dy)+y)
		d += fmt.Sprintf("z")
	} else {
		// ring drawn as inner radius -> outer radius -> outer arc -> inner radius -> inner arc
		rof, rif := float64(ro), float64(ri)
		dx, dy := p*math.Cos(w.Phi+math.Pi/2), p*math.Sin(w.Phi+math.Pi/2)
		a, b := int(rif*cphi+0.5+dx)+x, int(rif*sphi+0.5+dy)+y
		d = fmt.Sprintf("M %d,%d ", a, b)
		d += fmt.Sprintf("L %d,%d ", int(rof*cphi+0.5+dx)+x, int(rof*sphi+0.5+dy)+y)

		dx, dy = p*math.Cos(w.Psi-math.Pi/2), p*math.Sin(w.Psi-math.Pi/2)
		d += fmt.Sprintf("A %d,%d 0 0 1 %d,%d ", ro, ro, int(rof*cpsi+0.5+dx)+x, int(rof*spsi+0.5+dy)+y)
		d += fmt.Sprintf("L %d,%d ", int(rif*cpsi+0.5+dx)+x, int(rif*spsi+0.5+dy)+y)
		d += fmt.Sprintf("A %d,%d 0 0 0 %d,%d ", ri, ri, a, b)
		d += fmt.Sprintf("z")

	}

	sg.svg.Path(d, s+sf)
	 *************************/
}
