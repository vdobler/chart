package chart

import (
	"fmt"
	"math"
	"github.com/ajstarks/svgo"
)


// SvgGraphics implements BasicGraphics and uses the generic implementations
type SvgGraphics struct {
	svg  *svg.SVG
	w, h int
	font string
	fs   int
}

func NewSvgGraphics(sp *svg.SVG, width, height int, font string, fontsize int) *SvgGraphics {
	if font == "" {
		font = "Helvetica"
	}
	if fontsize == 0 {
		fontsize = 12
	}
	s := SvgGraphics{svg: sp, w: width, h: height, font: font, fs: fontsize}
	return &s
}

func (sg *SvgGraphics) Begin() {
	font, fs := sg.font, sg.fs
	if font == "" {
		font = "Arial"
	}
	if fs == 0 {
		fs = 12
	}
	sg.svg.Gstyle(fmt.Sprintf("stroke:#000000; stroke-width:1; font-family: %s; font-size: %d; opacity: 1; fill-opacity: 1", font, fs))
}

func (sg *SvgGraphics) End() {
	sg.svg.Gend()
}

func (sg *SvgGraphics) Dimensions() (int, int) {
	return sg.w, sg.h
}

func (sg *SvgGraphics) fontheight(font Font) (fh int) {
	if sg.fs <= 14 {
		fh = sg.fs + font.Size
	} else if sg.fs <= 20 {
		fh = sg.fs + 2*font.Size
	} else {
		fh = sg.fs + 3*font.Size
	}

	if fh == 0 {
		fh = 12
	}
	return
}

func (sg *SvgGraphics) FontMetrics(font Font) (fw float32, fh int, mono bool) {
	if font.Name == "" {
		font.Name = sg.font
	}
	fh = sg.fontheight(font)

	switch font.Name {
	case "Arial":
		fw, mono = 0.5*float32(fh), false
	case "Helvetica":
		fw, mono = 0.5*float32(fh), false
	case "Times":
		fw, mono = 0.51*float32(fh), false
	case "Courier":
		fw, mono = 0.62*float32(fh), true
	default:
		fw, mono = 0.75*float32(fh), false
	}

	// fmt.Printf("FontMetric of %s/%d: %.1f x %d  %t\n", style.Font, style.FontSize, fw, fh, mono)
	return
}

func (sg *SvgGraphics) TextLen(t string, font Font) int {
	return GenericTextLen(sg, t, font)
}


func (sg *SvgGraphics) Line(x0, y0, x1, y1 int, style DataStyle) {
	var s string
	if style.LineColor != "" {
		s = fmt.Sprintf("stroke:%s; ", style.LineColor)
	}
	s += fmt.Sprintf("stroke-width: %d; ", style.LineWidth)
	s += fmt.Sprintf("opacity: %.2f; ", 1-style.Alpha)

	sg.svg.Line(x0, y0, x1, y1, s)
}

func (sg *SvgGraphics) Text(x, y int, t string, align string, rot int, f Font) {
	if len(align) == 1 {
		align = "c" + align
	}
	_, fh, _ := sg.FontMetrics(f)

	trans := ""
	if rot != 0 {
		trans = fmt.Sprintf("transform=\"rotate(%d %d %d)\"", -rot, x, y)
	}

	// Hack because baseline alignments in svg often broken
	switch align[0] {
	case 'b':
		y += 0
	case 't':
		y += fh
	default:
		y += (4 * fh) / 10 // centered
	}
	s := "text-anchor:"
	switch align[1] {
	case 'l':
		s += "begin"
	case 'r':
		s += "end"
	default:
		s += "middle"
	}
	if f.Color != "" {
		s += "; stroke:" + f.Color
	}
	if f.Name != "" {
		s += "; font-family:" + f.Name
	}
	if f.Size != 0 {
		s += fmt.Sprintf("; font-size: %d", fh)
	}

	sg.svg.Text(x, y, t, trans, s)
}

func (sg *SvgGraphics) Symbol(x, y, s int, style DataStyle) {
	st := ""
	filled := "fill:solid"
	empty := "fill:none"
	if style.SymbolColor != "" {
		st += "stroke:" + style.SymbolColor
		filled = "fill:" + style.SymbolColor
	}
	f := style.SymbolSize
	if f == 0 {
		f = 1
	}
	lw := max(1, style.LineWidth)

	const n = 5               // default size
	a := int(n*f + 0.5)       // standard
	b := int(n/2*f + 0.5)     // smaller
	c := int(1.155*n*f + 0.5) // triangel long sist
	d := int(0.577*n*f + 0.5) // triangle short dist
	e := int(0.866*n*f + 0.5) // diagonal

	sg.svg.Gstyle(fmt.Sprintf("%s; stroke-width: %d", st, lw))
	switch style.Symbol {
	case '*':
		sg.svg.Line(x-e, y-e, x+e, y+e)
		sg.svg.Line(x-e, y+e, x+e, y-e)
		fallthrough
	case '+':
		sg.svg.Line(x-a, y, x+a, y)
		sg.svg.Line(x, y-a, x, y+a)
	case 'X':
		sg.svg.Line(x-e, y-e, x+e, y+e)
		sg.svg.Line(x-e, y+e, x+e, y-e)
	case 'o':
		sg.svg.Circle(x, y, a, empty)
	case '0':
		sg.svg.Circle(x, y, a, empty)
		sg.svg.Circle(x, y, b, empty)
	case '.':
		sg.svg.Circle(x, y, b, empty)
	case '@':
		sg.svg.Circle(x, y, a, filled)
	case '=':
		sg.svg.Rect(x-e, y-e, 2*e, 2*e, empty)
	case '#':
		sg.svg.Rect(x-e, y-e, 2*e, 2*e, filled)
	case 'A':
		sg.svg.Polygon([]int{x - a, x + a, x}, []int{y + d, y + d, y - c}, filled)
	case '%':
		sg.svg.Polygon([]int{x - a, x + a, x}, []int{y + d, y + d, y - c}, empty)
	case 'W':
		sg.svg.Polygon([]int{x - a, x + a, x}, []int{y - c, y - c, y + d}, filled)
	case 'V':
		sg.svg.Polygon([]int{x - a, x + a, x}, []int{y - c, y - c, y + d}, empty)
	case 'Z':
		sg.svg.Polygon([]int{x - e, x, x + e, x}, []int{y, y + e, y, y - e}, filled)
	case '&':
		sg.svg.Polygon([]int{x - e, x, x + e, x}, []int{y, y + e, y, y - e}, empty)
	default:
		sg.svg.Text(x, y, "?", "text-anchor:middle; alignment-baseline:middle")
	}
	sg.svg.Gend()

}

func (sg *SvgGraphics) Rect(x, y, w, h int, style DataStyle) {
	var s string
	linecol := style.LineColor
	if linecol != "" {
		s = fmt.Sprintf("stroke:%s; ", linecol)
	} else {
		linecol = "#808080"
	}
	s += fmt.Sprintf("stroke-width: %d; ", style.LineWidth)
	s += fmt.Sprintf("opacity: %.2f; ", 1-style.Alpha)
	if style.FillColor != "" {
		s += fmt.Sprintf("fill: %s; fill-opacity: %.2f", style.FillColor, 1-style.Alpha)
	} else {
		s += "fill-opacity: 0"
	}
	sg.svg.Rect(x, y, w, h, s)
	// GenericRect(sg, x, y, w, h, style) // TODO
}

func (sg *SvgGraphics) Style(element string) DataStyle {
	if v, ok := DefaultStyle[element]; ok {
		return v
	}
	return DataStyle{Symbol: 'o', SymbolColor: "#808080", LineColor: "#808080", LineWidth: 1, LineStyle: SolidLine}
}

func (sg *SvgGraphics) Font(element string) Font {
	if v, ok := DefaultFont[element]; ok {
		return v
	}
	return Font{}
}

func (sg *SvgGraphics) Title(text string) {
	_, fh, _ := sg.FontMetrics(Font{})
	x, y := sg.w/2, fh/2
	sg.Text(x, y, text, "tc", 0, sg.Font("title"))
}

func (sg *SvgGraphics) XAxis(xr Range, ys, yms int) {
	GenericXAxis(sg, xr, ys, yms)
}
func (sg *SvgGraphics) YAxis(yr Range, xs, xms int) {
	GenericYAxis(sg, yr, xs, xms)
}

func (sg *SvgGraphics) Scatter(points []EPoint, plotstyle PlotStyle, style DataStyle) {
	GenericScatter(sg, points, plotstyle, style)
}

func (sg *SvgGraphics) Boxes(boxes []Box, width int, style DataStyle) {
	GenericBoxes(sg, boxes, width, style)
}

func (sg *SvgGraphics) Key(x, y int, key Key) {
	GenericKey(sg, x, y, key)
}

func (sg *SvgGraphics) Bars(bars []Barinfo, style DataStyle) {
	GenericBars(sg, bars, style)
}

func (sg *SvgGraphics) Wedge(x, y, r int, phi, psi float64, style DataStyle) {
	d := fmt.Sprintf("M%d,%d ", x, y)
	rf := float64(r)
	d += fmt.Sprintf("L %d,%d", int(rf*math.Cos(phi)+0.5)+x, int(rf*math.Sin(phi)+0.5)+y)
	d += fmt.Sprintf("A %d,%d 0 0 1 %d,%d", r, r, int(rf*math.Cos(psi)+0.5)+x, int(rf*math.Sin(psi)+0.5)+y)
	d += fmt.Sprintf("z")

	var s string
	linecol := style.LineColor
	if linecol != "" {
		s = fmt.Sprintf("stroke:%s; ", linecol)
	} else {
		linecol = "#808080"
	}
	s += fmt.Sprintf("stroke-width: %d; ", style.LineWidth)
	s += fmt.Sprintf("opacity: %.2f; ", 1-style.Alpha)
	if style.FillColor != "" {
		s += fmt.Sprintf("fill: %s; fill-opacity: %.2f", style.FillColor, 1-style.Alpha)
	} else {
		s += "fill-opacity: 0"
	}

	sg.svg.Path(d, s)
	// GenericWedge(sg, x, y, r, phi, psi, style)
}
