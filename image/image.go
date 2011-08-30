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
func AddTo(img *image.RGBA, x, y, width, height int) *ImageGraphics {
	bg := img.At(x, y).(image.RGBAColor)
	return &ImageGraphics{Image: img, x0: x, y0: y, w: width, h: height, bg: bg}
}


func (ig *ImageGraphics) Begin() {
}

func (ig *ImageGraphics) End() {
}

func (ig *ImageGraphics) Dimensions() (int, int) {
	return ig.w, ig.h
}

func (ig *ImageGraphics) fontheight(font chart.Font) (fh int) {
	return 15
}

func (ig *ImageGraphics) FontMetrics(font chart.Font) (fw float32, fh int, mono bool) {
	return 8, 15, true
}

func (ig *ImageGraphics) TextLen(t string, font chart.Font) int {
	return chart.GenericTextLen(ig, t, font)
}


func (ig *ImageGraphics) Line(x0, y0, x1, y1 int, style chart.Style) {
	r, g, b := chart.Color2rgb(style.LineColor)
	col := image.RGBAColor{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(255 * style.Alpha)}

	// handle trivial cases first
	if x0 == x1 {
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		for ; y0 <= y1; y0++ {
			ig.Image.Set(ig.x0+x0, ig.y0+y0, col)
		}
		return
	}
	if y0 == y1 {
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		for ; x0 <= x1; x0++ {
			ig.Image.Set(ig.x0+x0, ig.y0+y0, col)
		}
		return
	}
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := sign(x1-x0), sign(y1-y0)
	err, e2 := dx+dy, 0
	for {
		ig.Image.Set(ig.x0+x0, ig.y0+y0, col)
		// fmt.Printf("%d %d   %d %d\n", x0,y0, x1, y1)
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

func (ig *ImageGraphics) Text(x, y int, t string, align string, rot int, f chart.Font) {
	if len(align) == 1 {
		align = "c" + align
	}
	fw, fh, _ := ig.FontMetrics(f)
	//fmt.Printf("Text '%s' at (%d,%d) %s\n", t, x,y, align)
	// TODO: handle rot

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

	color := "#000000"
	if f.Color != "" {
		color = f.Color
	}
	color += ""
	var R, G, B uint32
	R, G, B = 0, 0, 0 // TODO read from font color

	// ig.Text(x, y, t, trans, s)
	for i, c := range t {
		if _, ok := font[c]; !ok {
			c = '?'
		}
		for l := 0; l < 15; l++ {
			q := font[c][l]
			for bit := 7; bit >= 0; bit-- {
				v := uint32(q & 0xff)
				if v > 0 {
					xx, yy := ig.x0+x+i*8+bit, ig.y0+y+l
					r, g, b, _ := ig.Image.At(xx, yy).RGBA()
					r >>= 8
					g >>= 8
					b >>= 8
					r *= 0xff - v
					g *= 0xff - v
					b *= 0xff - v
					r += R * v
					g += G * v
					b += B * v
					r >>= 8
					g >>= 8
					b >>= 8
					ig.Image.Set(xx, yy, image.RGBAColor{uint8(r), uint8(g), uint8(b), 0})
				}
				q >>= 8
			}
		}
	}

}

func (ig *ImageGraphics) Symbol(x, y int, style chart.Style) {
	chart.GenericSymbol(ig, x, y, style)
}

func (ig *ImageGraphics) Rect(x, y, w, h int, style chart.Style) {
	chart.GenericRect(ig, x, y, w, h, style)
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
