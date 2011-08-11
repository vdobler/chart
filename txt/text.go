package chart

import (
	"fmt"
	"math"
	"github.com/vdobler/chart"
)

// Different edge styles for boxes
var Edge = [][4]int{{'+', '+', '+', '+'}, {'.', '.', '\'', '\''}, {'/', '\\', '\\', '/'}}


// A Text Buffer
type TextBuf struct {
	Buf  []int // the internal buffer.  Point (x,y) is mapped to x + y*(W+1)
	W, H int   // Width and Height
}

// Set up a new TextBuf with width w and height h.
func NewTextBuf(w, h int) (tb *TextBuf) {
	tb = new(TextBuf)
	tb.W, tb.H = w, h
	tb.Buf = make([]int, (w+1)*h)
	for i, _ := range tb.Buf {
		tb.Buf[i] = ' '
	}
	for i := 0; i < h; i++ {
		tb.Buf[i*(w+1)+w] = '\n'
	}
	// tb.Buf[0], tb.Buf[(w+1)*h-1] = 'X', 'X'
	return
}


// Put character c at (x,y)
func (tb *TextBuf) Put(x, y, c int) {
	if x < 0 || y < 0 || x >= tb.W || y >= tb.H {
		return
		// fmt.Printf("Ooooops Put(): %d, %d  '%c' \n", x, y, c)
		x, y = 0, 0
	}
	i := (tb.W+1)*y + x
	tb.Buf[i] = c
}

// Draw rectangle of width w and height h from corner at (x,y).
// Use one of the corner style defined in Edge. 
// Interior is filled with charater fill iff != 0.
func (tb *TextBuf) Rect(x, y, w, h int, style int, fill int) {
	style = style % len(Edge)

	if h < 0 {
		h = -h
		y -= h
	}
	if w < 0 {
		w = -w
		x -= w
	}

	tb.Put(x, y, Edge[style][0])
	tb.Put(x+w, y, Edge[style][1])
	tb.Put(x, y+h, Edge[style][2])
	tb.Put(x+w, y+h, Edge[style][3])
	for i := 1; i < w; i++ {
		tb.Put(x+i, y, '-')
		tb.Put(x+i, y+h, '-')
	}
	for i := 1; i < h; i++ {
		tb.Put(x+w, y+i, '|')
		tb.Put(x, y+i, '|')
		if fill > 0 {
			for j := x + 1; j < x+w; j++ {
				tb.Put(j, y+i, fill)
			}
		}
	}
}

func (tb *TextBuf) Block(x, y, w, h int, fill int) {
	if h < 0 {
		h = -h
		y -= h
	}
	if w < 0 {
		w = -w
		x -= w
	}
	for i := 0; i < w; i++ {
		for j := 0; j <= h; j++ {
			tb.Put(x+i, y+j, fill)
		}
	}
}

// Return real character len of s (rune count).
func StrLen(s string) (n int) {
	for _, _ = range s {
		n++
	}
	return
}

// Print text txt at (x,y). Horizontal display for align in [-1,1],
// vasrtical alignment for align in [2,4]
// align: -1: left; 0: centered; 1: right; 2: top, 3: center, 4: bottom
func (tb *TextBuf) Text(x, y int, txt string, align int) {
	if align <= 1 {
		switch align {
		case 0:
			x -= StrLen(txt) / 2
		case 1:
			x -= StrLen(txt)
		}
		i := 0
		for _, rune := range txt {
			tb.Put(x+i, y, rune)
			i++
		}
	} else {
		switch align {
		case 3:
			y -= StrLen(txt) / 2
		case 2:
			x -= StrLen(txt)
		}
		i := 0
		for _, rune := range txt {
			tb.Put(x, y+i, rune)
			i++
		}
	}
}


// Paste buf at (x,y)
func (tb *TextBuf) Paste(x, y int, buf *TextBuf) {
	s := buf.W + 1
	for i := 0; i < buf.W; i++ {
		for j := 0; j < buf.H; j++ {
			tb.Put(x+i, y+j, buf.Buf[i+s*j])
		}
	}
}

func (tb *TextBuf) Line(x0, y0, x1, y1 int, symbol int) {
	// handle trivial cases first
	if x0 == x1 {
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		for ; y0 <= y1; y0++ {
			tb.Put(x0, y0, symbol)
		}
		return
	}
	if y0 == y1 {
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		for ; x0 <= x1; x0++ {
			tb.Put(x0, y0, symbol)
		}
		return
	}
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := sign(x1-x0), sign(y1-y0)
	err, e2 := dx+dy, 0
	for {
		tb.Put(x0, y0, symbol)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 = 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}

	}
}


// Convert to plain (utf8) string.
func (tb *TextBuf) String() string {
	return string(tb.Buf)
}


// TextGraphics
type TextGraphics struct {
	tb   *TextBuf
	w, h int
}

func NewTextGraphics(w, h int) *TextGraphics {
	tg := TextGraphics{}
	tg.tb = NewTextBuf(w, h)
	tg.w, tg.h = w, h
	return &tg
}


func (g *TextGraphics) Begin() {
	g.tb = NewTextBuf(g.w, g.h)
}

func (g *TextGraphics) End() {}
func (g *TextGraphics) Dimensions() (int, int) {
	return g.w, g.h
}
func (g *TextGraphics) FontMetrics(font chart.Font) (fw float32, fh int, mono bool) {
	return 1, 1, true
}

func (g *TextGraphics) TextLen(t string, font chart.Font) int {
	return len(t)
}


func (g *TextGraphics) Line(x0, y0, x1, y1 int, style chart.Style) {
	symbol := style.Symbol
	if symbol < ' ' || symbol > '~' {
		symbol = 'x'
	}
	g.tb.Line(x0, y0, x1, y1, symbol)
}

func (g *TextGraphics) Text(x, y int, t string, align string, rot int, font chart.Font) {
	// align: -1: left; 0: centered; 1: right; 2: top, 3: center, 4: bottom
	if len(align) == 2 {
		align = align[1:]
	}
	a := 0
	if rot == 0 {
		if align == "l" {
			a = -1
		}
		if align == "c" {
			a = 0
		}
		if align == "r" {
			a = 1
		}
	} else {
		if align == "l" {
			a = 2
		}
		if align == "c" {
			a = 3
		}
		if align == "r" {
			a = 4
		}
	}
	g.tb.Text(x, y, t, a)
}

func (g *TextGraphics) Rect(x, y, w, h int, style chart.Style) {
	// Normalize coordinates
	if h < 0 {
		h = -h
		y -= h
	}
	if w < 0 {
		w = -w
		x -= w
	}

	// Border
	if style.LineWidth > 0 {
		for i := 0; i < w; i++ {
			g.tb.Put(x+i, y, style.Symbol)
			g.tb.Put(x+i, y+h-1, style.Symbol)
		}
		for i := 1; i < h-1; i++ {
			g.tb.Put(x, y+i, style.Symbol)
			g.tb.Put(x+w-1, y+i, style.Symbol)
		}
	}

	// Filling
	if style.FillColor != "" {
		// TODO: fancier logic
		var s int
		if style.FillColor == "#000000" {
			s = '#' // black
		} else if style.FillColor == "#ffffff" {
			s = ' ' // white
		} else {
			s = style.Symbol
		}
		for i := 1; i < h-1; i++ {
			for j := 1; j < w-1; j++ {
				g.tb.Put(x+j, y+i, s)
			}
		}
	}
}

func (g *TextGraphics) Style(element string) chart.Style {
	b := "#000000"
	return chart.Style{Symbol: 'o', SymbolColor: b, LineColor: b, LineWidth: 1, LineStyle: chart.SolidLine}
}

func (g *TextGraphics) Font(element string) chart.Font {
	return chart.Font{}
}

func (g *TextGraphics) String() string {
	return g.tb.String()
}

func (g *TextGraphics) Symbol(x, y, s int, style chart.Style) {
	g.tb.Put(x, y, s)
}
func (g *TextGraphics) Title(text string) {
	x, y := g.w/2, 1
	g.Text(x, y, text, "tc", 0, chart.Font{})
}

func (g *TextGraphics) XAxis(xrange chart.Range, y, y1 int) {
	mirror := xrange.TicSetting.Mirror
	xa, xe := xrange.Data2Screen(xrange.Min), xrange.Data2Screen(xrange.Max)
	for sx := xa; sx <= xe; sx++ {
		g.tb.Put(sx, y, '-')
		if mirror >= 1 {
			g.tb.Put(sx, y1, '-')
		}
	}
	if xrange.ShowZero && xrange.Min < 0 && xrange.Max > 0 {
		z := xrange.Data2Screen(0)
		for yy := y - 1; yy > y1+1; yy-- {
			g.tb.Put(z, yy, ':')
		}
	}

	if xrange.Label != "" {
		yy := y + 1
		if !xrange.TicSetting.Hide {
			yy++
		}
		g.tb.Text((xa+xe)/2, yy, xrange.Label, 0)
	}

	for _, tic := range xrange.Tics {
		x := xrange.Data2Screen(tic.Pos)
		lx := xrange.Data2Screen(tic.LabelPos)
		if xrange.Time {
			g.tb.Put(x, y, '|')
			if mirror >= 2 {
				g.tb.Put(x, y1, '|')
			}
			g.tb.Put(x, y+1, '|')
			if tic.Align == -1 {
				g.tb.Text(lx+1, y+1, tic.Label, -1)
			} else {
				g.tb.Text(lx, y+1, tic.Label, 0)
			}
		} else {
			g.tb.Put(x, y, '+')
			if mirror >= 2 {
				g.tb.Put(x, y1, '+')
			}
			g.tb.Text(lx, y+1, tic.Label, 0)
		}
		if xrange.ShowLimits {
			if xrange.Time {
				g.tb.Text(xa, y+2, xrange.TMin.Format("2006-01-02 15:04:05"), -1)
				g.tb.Text(xe, y+2, xrange.TMax.Format("2006-01-02 15:04:05"), 1)
			} else {
				g.tb.Text(xa, y+2, fmt.Sprintf("%g", xrange.Min), -1)
				g.tb.Text(xe, y+2, fmt.Sprintf("%g", xrange.Max), 1)
			}
		}
	}

	// GenericXAxis(g, xr, ys, yms)
}
func (g *TextGraphics) YAxis(yrange chart.Range, x, x1 int) {
	label := yrange.Label
	mirror := yrange.TicSetting.Mirror
	ya, ye := yrange.Data2Screen(yrange.Min), yrange.Data2Screen(yrange.Max)
	for sy := min(ya, ye); sy <= max(ya, ye); sy++ {
		g.tb.Put(x, sy, '|')
		if mirror >= 1 {
			g.tb.Put(x1, sy, '|')
		}
	}
	if yrange.ShowZero && yrange.Min < 0 && yrange.Max > 0 {
		z := yrange.Data2Screen(0)
		for xx := x + 1; xx < x1; xx += 2 {
			g.tb.Put(xx, z, '-')
		}
	}

	if label != "" {
		g.tb.Text(1, (ya+ye)/2, label, 3)
	}

	for _, tic := range yrange.Tics {
		y := yrange.Data2Screen(tic.Pos)
		ly := yrange.Data2Screen(tic.LabelPos)
		if yrange.Time {
			g.tb.Put(x, y, '+')
			if mirror >= 2 {
				g.tb.Put(x1, y, '+')
			}
			if tic.Align == 0 { // centered tic
				g.tb.Put(x-1, y, '-')
				g.tb.Put(x-2, y, '-')
			}
			g.tb.Text(x, ly, tic.Label+" ", 1)
		} else {
			g.tb.Put(x, y, '+')
			if mirror >= 2 {
				g.tb.Put(x1, y, '+')
			}
			g.tb.Text(x-2, ly, tic.Label, 1)
		}
	}
}

func (g *TextGraphics) Scatter(points []chart.EPoint, plotstyle chart.PlotStyle, style chart.Style) {
	// First pass: Error bars
	for _, p := range points {
		xl, yl, xh, yh := p.BoundingBox()
		if !math.IsNaN(p.DeltaX) {
			g.tb.Line(int(xl), int(p.Y), int(xh), int(p.Y), '-')
		}
		if !math.IsNaN(p.DeltaY) {
			g.tb.Line(int(p.X), int(yl), int(p.X), int(yh), '|')
		}
	}

	// Second pass: Line
	if (plotstyle&chart.PlotStyleLines) != 0 && len(points) > 0 {
		lastx, lasty := points[0].X, points[0].Y
		for i := 1; i < len(points); i++ {
			x, y := points[i].X, points[i].Y
			g.tb.Line(int(lastx), int(lasty), int(x), int(y), style.Symbol)
			lastx, lasty = x, y
		}
	}

	// Third pass: symbols
	if (plotstyle&chart.PlotStylePoints) != 0 && len(points) != 0 {
		for _, p := range points {
			g.tb.Put(int(p.X), int(p.Y), style.Symbol)
		}
	}
	// chart.GenericScatter(g, points, plotstyle, style)
}

func (g *TextGraphics) Boxes(boxes []chart.Box, width int, style chart.Style) {
	if width%2 == 0 {
		width += 1
	}
	hbw := (width - 1) / 2
	if style.Symbol == 0 {
		style.Symbol = '*'
	}

	for _, box := range boxes {
		x := int(box.X)
		q1, q3 := int(box.Q1), int(box.Q3)
		g.tb.Rect(x-hbw, q1, 2*hbw, q3-q1, 0, ' ')
		if !math.IsNaN(box.Med) {
			med := int(box.Med)
			g.tb.Put(x-hbw, med, '+')
			for i := 0; i < hbw; i++ {
				g.tb.Put(x-i, med, '-')
				g.tb.Put(x+i, med, '-')
			}
			g.tb.Put(x+hbw, med, '+')
		}

		if !math.IsNaN(box.Avg) && style.Symbol != 0 {
			g.tb.Put(x, int(box.Avg), style.Symbol)
		}

		if !math.IsNaN(box.High) {
			for y := int(box.High); y < q3; y++ {
				g.tb.Put(x, y, '|')
			}
		}

		if !math.IsNaN(box.Low) {
			for y := int(box.Low); y > q1; y-- {
				g.tb.Put(x, y, '|')
			}
		}

		for _, ol := range box.Outliers {
			y := int(ol)
			g.tb.Put(x, y, style.Symbol)
		}
	}
}

var (
	KeyHorSep      float32 = 1.5
	KeyVertSep     float32 = 0.5
	KeyColSep      float32 = 2.0
	KeySymbolWidth float32 = 4
	KeySymbolSep   float32 = 1
	KeyRowSep      float32 = 0.75
)


func (g *TextGraphics) Key(x, y int, key chart.Key) {
	m := key.Place()
	tw, th, cw, rh := key.Layout(g, m)
	style := g.Style("key")
	if style.LineWidth > 0 || style.FillColor != "" {
		g.tb.Rect(x, y, tw, th, 1, ' ')
	}
	x += int(KeyHorSep)
	vsep := KeyVertSep
	if vsep < 1 {
		vsep = 1
	}
	y += int(vsep)
	for ci, col := range m {
		yy := y

		for ri, e := range col {
			if e == nil || e.Text == "" {
				continue
			}
			plotStyle := e.PlotStyle
			// fmt.Printf("KeyEntry %s: PlotStyle = %d\n", e.Text, e.PlotStyle)
			if plotStyle == -1 {
				// heading only...
				g.tb.Text(x, yy, e.Text, -1)
			} else {
				// normal entry
				if (plotStyle & chart.PlotStyleLines) != 0 {
					g.Line(x, yy, x+int(KeySymbolWidth), yy, e.Style)
				}
				if (plotStyle & chart.PlotStylePoints) != 0 {
					g.Symbol(x+int(KeySymbolWidth/2), yy, e.Style.Symbol, e.Style)
				}
				if (plotStyle & chart.PlotStyleBox) != 0 {
					g.tb.Put(x+int(KeySymbolWidth/2), yy, e.Style.Symbol)
				}
				g.tb.Text(x+int((KeySymbolWidth+KeySymbolSep)), yy, e.Text, -1)
			}
			yy += rh[ri] + int(KeyRowSep)
		}

		x += int((KeySymbolWidth + KeySymbolSep + KeyColSep + float32(cw[ci])))
	}

}

func (g *TextGraphics) Bars(bars []chart.Barinfo, style chart.Style) {
	chart.GenericBars(g, bars, style)
}

func (g *TextGraphics) Wedge(x, y, ry int, phi, psi float64, style chart.Style) {
	rx := int(1.9 * float64(ry))
	x += 10 // TODO: find a proper way here....
	ryf, rxf := float64(ry), float64(rx)
	xa, ya := int(math.Cos(phi)*rxf)+x, int(math.Sin(phi)*ryf)+y
	xc, yc := int(math.Cos(psi)*rxf)+x, int(math.Sin(psi)*ryf)+y

	if math.Fabs(phi-psi) >= 4*math.Pi {
		phi, psi = 0, 2*math.Pi
	} else {
		g.Line(x, y, xa, ya, style)
		g.Line(x, y, xc, yc, style)
	}

	if style.FillColor != "" {
		delta := 1 / (4 * rxf)
		ls := chart.Style{LineColor: style.FillColor, LineWidth: 1, Symbol: style.Symbol}
		for a := phi; a <= psi; a += delta {
			xr, yr := int(math.Cos(a)*rxf)+x, int(math.Sin(a)*ryf)+y
			g.Line(x, y, xr, yr, ls)
		}
	}

	var xb, yb int
	for ; phi < psi; phi += 0.1 { // aproximate circle by 62-corner
		xb, yb = int(math.Cos(phi)*rxf)+x, int(math.Sin(phi)*ryf)+y
		g.Line(xa, ya, xb, yb, style)
		xa, ya = xb, yb
	}
	g.Line(xb, yb, xc, yc, style)

}

func (g *TextGraphics) Rings(wedges []chart.Wedgeinfo, x, y, r int) {
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
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

func sign(a int) int {
	if a < 0 {
		return -1
	}
	if a == 0 {
		return 0
	}
	return 1
}
