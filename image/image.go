package imgg

import (
	"image"
	"fmt"

	"github.com/vdobler/chart"
)

// ImageGraphics implements BasicGraphics and uses the generic implementations
type ImageGraphics struct {
	Image  *image.RGBA
	x0, y0 int
	w, h   int
	bg     image.RGBAColor
}

// New creates a new ImageGraphics of dimension w x h.
func New(width, height int, background image.RGBAColor) *ImageGraphics {
	_ = fmt.Sprintf("")
	img := image.NewRGBA(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, background)
		}
	}
	return &ImageGraphics{Image: img, x0: 0, y0: 0, w: width, h: height, bg: background}
}

// AddTo returns a new ImageGraphics which will write to (width x height) sized
// area starting at (x,y) on the provided image img.
func AddTo(img *image.RGBA, x, y, width, height int, background image.RGBAColor) *ImageGraphics {
	return &ImageGraphics{Image: img, x0: x, y0: y, w: width, h: height, bg: background}
}

func (ig *ImageGraphics) Begin()                              {}
func (ig *ImageGraphics) End()                                {}
func (ig *ImageGraphics) Background() (r, g, b, a uint8)      { return ig.bg.R, ig.bg.G, ig.bg.B, ig.bg.A }
func (ig *ImageGraphics) Dimensions() (int, int)              { return ig.w, ig.h }
func (ig *ImageGraphics) fontheight(font chart.Font) (fh int) { return 15 }
func (ig *ImageGraphics) FontMetrics(font chart.Font) (fw float32, fh int, mono bool) {
	return 8, 15, true
}
func (ig *ImageGraphics) TextLen(t string, font chart.Font) int {
	return chart.GenericTextLen(ig, t, font)
}

func ddAndPat(style chart.Style) (d, dd int, pat []bool) {
	var ok bool
	if pat, ok = dashPattern[style.LineStyle]; !ok {
		pat = dashPattern[chart.SolidLine]
	}

	d = (style.LineWidth - 1) / 2
	dd = d
	if style.LineWidth%2 == 0 {
		dd++
	}
	return
}

func (ig *ImageGraphics) Line(x0, y0, x1, y1 int, style chart.Style) {
	d, dd, pat := ddAndPat(style)
	for xd := -d; xd <= dd; xd++ {
		for yd := -d; yd <= dd; yd++ {
			ig.oneLine(x0+xd, y0+yd, x1+xd, y1+yd, style, pat, 0)
		}
	}
}

var dashPattern map[int][]bool = map[int][]bool{
	chart.SolidLine:      []bool{true},
	chart.DashedLine:     []bool{true, true, true, true, true, true, true, true, false, false, false},
	chart.DottedLine:     []bool{true, true, true, false, false, false},
	chart.DashDotDotLine: []bool{true, true, true, true, true, true, false, false, true, true, false, false},
	chart.LongDashLine: []bool{true, true, true, true, true, true, true, true,
		false, false, false, false, false, false, false, false},
	chart.LongDotLine: []bool{true, true, true, false, false, false, false, false, false},
}

func (ig *ImageGraphics) oneLine(x0, y0, x1, y1 int, style chart.Style, pat []bool, P int) int {
	R, G, B := chart.Color2rgb(style.LineColor)
	r, g, b := uint32(R), uint32(G), uint32(B)
	alpha := uint32(0xff * style.Alpha)
	N := len(pat)

	// handle trivial cases first
	if x0 == x1 {
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		for ; y0 <= y1; y0++ {
			if pat[P%N] {
				ig.paint(ig.x0+x0, ig.y0+y0, r, g, b, alpha)
				// ig.Image.Set(ig.x0+x0, ig.y0+y0, col)
			}
			P++
		}
		return P
	}
	if y0 == y1 {
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		for ; x0 <= x1; x0++ {
			if pat[P%N] {
				ig.paint(ig.x0+x0, ig.y0+y0, r, g, b, alpha)
				// ig.Image.Set(ig.x0+x0, ig.y0+y0, col)
			}
			P++
		}
		return P
	}
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := sign(x1-x0), sign(y1-y0)
	err, e2 := dx+dy, 0
	for {
		if pat[P%N] {
			ig.paint(ig.x0+x0, ig.y0+y0, r, g, b, alpha)
			// ig.Image.Set(ig.x0+x0, ig.y0+y0, col)
		}
		P++
		// fmt.Printf("%d %d   %d %d\n", x0,y0, x1, y1)
		if x0 == x1 && y0 == y1 {
			return P
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
	return 0
}

func (ig *ImageGraphics) Path(x, y []int, style chart.Style) {
	n := min(len(x), len(y))
	d, dd, pat := ddAndPat(style)

	p := 0
	for i := 1; i < n; i++ {
		x0, y0 := x[i-1], y[i-1]
		x1, y1 := x[i], y[i]
		np := 0
		for xd := -d; xd <= dd; xd++ {
			for yd := -d; yd <= dd; yd++ {
				np = ig.oneLine(x0+xd, y0+yd, x1+xd, y1+yd, style, pat, p)
			}
		}
		p = np
	}
}

func (ig *ImageGraphics) Text(x, y int, t string, align string, rot int, f chart.Font) {
	if len(align) == 1 {
		align = "c" + align
	}
	fw, fh, _ := ig.FontMetrics(f)
	//fmt.Printf("Text '%s' at (%d,%d) %s\n", t, x,y, align)
	// TODO: handle rot

	if rot == 90 {
		switch align[0] {
		case 'b':
			x -= fh
		case 't':
			y += 0
		default:
			x -= fh / 2
		}
		// TODO: rune count
		switch align[1] {
		case 'l':
			y += 0
		case 'r':
			y += int(fw * float32(len(t)))
		default:
			y += int(fw / 2 * float32(len(t)))
		}
	} else {
		switch align[0] {
		case 'b':
			y -= fh
		case 't':
			y += 0
		default:
			y -= fh / 2
		}
		// TODO: rune count
		switch align[1] {
		case 'l':
			x += 0
		case 'r':
			x -= int(fw * float32(len(t)))
		default:
			x -= int(fw / 2 * float32(len(t)))
		}
	}

	color := "#000000"
	if f.Color != "" {
		color = f.Color
	}
	RR, GG, BB := chart.Color2rgb(color)
	R, G, B := uint32(RR), uint32(GG), uint32(BB)

	// ig.Text(x, y, t, trans, s)
	var xx, yy int
	for i, c := range t {
		if _, ok := font[c]; !ok {
			c = '?'
		}
		for l := 0; l < 15; l++ {
			q := font[c][l]
			for bit := 7; bit >= 0; bit-- {
				v := uint32(q & 0xff)
				if v > 0 {
					if rot == 90 {
						xx, yy = ig.x0+x+l, ig.y0+y-i*8-bit
					} else {
						xx, yy = ig.x0+x+i*8+bit, ig.y0+y+l
					}
					ig.paint(xx, yy, R, G, B, 0xff-v)
				}
				q >>= 8
			}
		}
	}

}

func (ig *ImageGraphics) paint(x, y int, R, G, B uint32, alpha uint32) {
	r, g, b, a := ig.Image.At(x, y).RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	r *= alpha
	g *= alpha
	b *= alpha
	a *= alpha
	r += R * (0xff - alpha)
	g += G * (0xff - alpha)
	b += B * (0xff - alpha)
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	ig.Image.Set(x, y, image.RGBAColor{uint8(r), uint8(g), uint8(b), uint8(a)})
}

func (ig *ImageGraphics) Symbol(x, y int, style chart.Style) {
	chart.GenericSymbol(ig, x, y, style)
}

func (ig *ImageGraphics) Rect(x, y, w, h int, style chart.Style) {
	chart.GenericRect(ig, x, y, w, h, style)
}

func (ig *ImageGraphics) Wedge(x, y, ro, ri int, phi, psi float64, style chart.Style) {
	chart.GenericWedge(ig, x, y, ro, ri, phi, psi, 1, style)
}

func (ig *ImageGraphics) Title(text string) {
	font := chart.DefaultFont["title"]
	_, fh, _ := ig.FontMetrics(font)
	x, y := ig.w/2, fh/2
	ig.Text(x, y, text, "tc", 0, font)
}

func (ig *ImageGraphics) XAxis(xr chart.Range, ys, yms int) {
	chart.GenericXAxis(ig, xr, ys, yms)
}
func (ig *ImageGraphics) YAxis(yr chart.Range, xs, xms int) {
	chart.GenericYAxis(ig, yr, xs, xms)
}

func (ig *ImageGraphics) Scatter(points []chart.EPoint, plotstyle chart.PlotStyle, style chart.Style) {
	chart.GenericScatter(ig, points, plotstyle, style)
}

func (ig *ImageGraphics) Boxes(boxes []chart.Box, width int, style chart.Style) {
	chart.GenericBoxes(ig, boxes, width, style)
}

func (ig *ImageGraphics) Key(x, y int, key chart.Key) {
	chart.GenericKey(ig, x, y, key)
}

func (ig *ImageGraphics) Bars(bars []chart.Barinfo, style chart.Style) {
	chart.GenericBars(ig, bars, style)
}

func (ig *ImageGraphics) Rings(wedges []chart.Wedgeinfo, x, y, ro, ri int) {
	chart.GenericRings(ig, wedges, x, y, ro, ri, 1)
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
