package chart

import (
	"fmt"
	"math"
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


func LayoutTxt(w, h int, title, xlabel, ylabel string, hidextics, hideytics bool, key *Key, fw, fh int) (width, leftm, height, topm int, kb *TextBuf, numxtics, numytics int) {
	if key.Pos == "" {
		key.Pos = "itr"
	}

	if h < 5 {
		h = 5
	}
	if w < 10 {
		w = 10
	}

	width, leftm, height, topm = w-6*fw, 2*fw, h-1*fh, 0
	xlabsep, ylabsep := fh, 3*fw
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
		leftm += 6 * fw
		width -= 6 * fw
		ylabsep += 6 * fw
	}

	if !key.Hide { // TODO: buggy, not device independent
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
	case width/fw < 20:
		numxtics = 2
	case width/fw < 30:
		numxtics = 3
	case width/fw < 60:
		numxtics = 4
	case width/fw < 80:
		numxtics = 5
	case width/fw < 100:
		numxtics = 7
	default:
		numxtics = 10
	}
	// fmt.Printf("Requesting %d,%d tics.\n", ntics,height/3)

	numytics = (h / fh) / 5

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
func (g *TextGraphics) FontMetrics(font Font) (fw float32, fh int, mono bool) {
	return 1, 1, true
}

func (g *TextGraphics) TextLen(t string, font Font) int {
	return len(t)
}


func (g *TextGraphics) Line(x0, y0, x1, y1 int, style DataStyle) {
	g.tb.Line(x0, y0, x1, y1, '*')
}

func (g *TextGraphics) Text(x, y int, t string, align string, rot int, font Font) {
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

func (g *TextGraphics) Rect(x, y, w, h int, style DataStyle) {
	g.tb.Rect(x, y, w, h, 1, ' ') // TODO use info from style
}

func (g *TextGraphics) Style(element string) DataStyle {
	b := "#000000"
	return DataStyle{Symbol: 'o', SymbolColor: b, LineColor: b, LineWidth: 1, LineStyle: SolidLine}
}

func (g *TextGraphics) Font(element string) Font {
	return Font{}
}

func (g *TextGraphics) String() string {
	return g.tb.String()
}

func (g *TextGraphics) Symbol(x, y, s int, style DataStyle) {
	g.tb.Put(x, y, s)
}
func (g *TextGraphics) Title(text string) {
	x, y := g.w/2, 0
	g.Text(x, y, text, "tc", 0, Font{})
}

func (g *TextGraphics) XAxis(xrange Range, y, y1 int) {
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
func (g *TextGraphics) YAxis(yrange Range, x, x1 int) {
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

func (g *TextGraphics) Scatter(points []EPoint, plotstyle PlotStyle, style DataStyle) {
	GenericScatter(g, points, plotstyle, style)
}

func (g *TextGraphics) Boxes(boxes []Box, width int, style DataStyle) {
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

func (g *TextGraphics) Key(x, y int, key Key) {
	GenericKey(g, x, y, key)
}

func (g *TextGraphics) Bars(bars []Barinfo, style DataStyle) {
	GenericBars(g, bars, style)
}

func (g *TextGraphics) Wedge(x, y, r int, phi, psi float64, style DataStyle) {
	GenericWedge(g, x, y, r, phi, psi, style)
}
