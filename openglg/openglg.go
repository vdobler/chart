package openglg

import (
	"image/color"
	"log"
	"runtime"

	"github.com/pwaller/go-chart"

	"github.com/banthar/gl"

	glh "github.com/pwaller/go-glhelpers"
)

var thickwhite = chart.Style{LineColor: "#ffffff", LineWidth: 2, Alpha: 1}
var thinwhite = chart.Style{LineColor: "#ffffff", LineWidth: 1, Alpha: 1}
var dark = chart.Style{FillColor: "#222222", LineWidth: 0, Alpha: 0.75, Font: chart.Font{Size: chart.HugeFontSize}}

var DarkStyle = map[chart.PlotElement]chart.Style{
	chart.MajorAxisElement: thickwhite,
	chart.MinorAxisElement: thinwhite,
	chart.MajorTicElement:  thickwhite,
	chart.MinorTicElement:  thinwhite,
	chart.ZeroAxisElement:  thinwhite,
	chart.GridLineElement:  thinwhite,
	chart.GridBlockElement: dark,
	chart.KeyElement:       dark,
}

// OpenGLGraphics implements BasicGraphics and uses the generic implementations
type OpenGLGraphics struct {
	w, h   int
	font   string
	fs     int
	bg     color.RGBA
	tx, ty int
}

func whoami() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		return "unnamed"
	}
	return me.Name()
}

// New creates a new OpenGLGraphics of dimension w x h, with a default font font of size fontsize.
func New(width, height int, font string, fontsize int, background color.RGBA) *OpenGLGraphics {
	if font == "" {
		font = "Helvetica"
	}
	if fontsize == 0 {
		fontsize = 12
	}
	s := OpenGLGraphics{w: width, h: height, font: font, fs: fontsize, bg: background}
	return &s
}

// AddTo returns a new ImageGraphics which will write to (width x height) sized
// area starting at (x,y) on the provided SVG
func AddTo(x, y, width, height int, font string, fontsize int, background color.RGBA) *OpenGLGraphics {
	log.Panicf("Unimplemented: %s", whoami())
	// s := New(sp, width, height, font, fontsize, background)
	// s.tx, s.ty = x, y
	// return s
	return nil
}

func (sg *OpenGLGraphics) Options() chart.PlotOptions {
	//log.Panicf("Unimplemented: %s", whoami())
	return nil
}

func (sg *OpenGLGraphics) Begin() {
	// TODO: background rect?
	return
	log.Panicf("Unimplemented: %s", whoami())
	// font, fs := sg.font, sg.fs
	// if font == "" {
	// 	font = "Helvetica"
	// }
	// if fs == 0 {
	// 	fs = 12
	// }
	// sg.svg.Gstyle(fmt.Sprintf("font-family: %s; font-size: %d",
	// 	font, fs))
	// if sg.tx != 0 || sg.ty != 0 {
	// 	sg.svg.Gtransform(fmt.Sprintf("translate(%d %d)", sg.tx, sg.ty))
	// }

	// bgc := fmt.Sprintf("#%02x%02x%02x", sg.bg.R>>8, sg.bg.G>>8, sg.bg.B>>8)
	// opa := fmt.Sprintf("%.4f", float64(sg.bg.A>>8)/255)
	// bgs := fmt.Sprintf("stroke: %s; opacity: %s; fill: %s; fill-opacity: %s", bgc, opa, bgc, opa)
	// sg.svg.Rect(0, 0, sg.w, sg.h, bgs)
}

func (sg *OpenGLGraphics) End() {
	// TODO: Anything?
	//log.Panicf("Unimplemented: %s", whoami())
	// sg.svg.Gend()
	// if sg.tx != 0 || sg.ty != 0 {
	// 	sg.svg.Gend()
	// }
}

func (sg *OpenGLGraphics) Background() (r, g, b, a uint8) {
	log.Panicf("Unimplemented: %s", whoami())
	// return uint8(sg.bg.R >> 8), uint8(sg.bg.G >> 8), uint8(sg.bg.B >> 8), uint8(sg.bg.A >> 8)
	return 0, 0, 0, 0
}

func (sg *OpenGLGraphics) Dimensions() (int, int) {
	//log.Panicf("Unimplemented: %s", whoami())
	return sg.w, sg.h
	//return 0, 0
}

func (sg *OpenGLGraphics) fontheight(font chart.Font) (fh int) {
	// log.Panicf("Unimplemented: %s", whoami())
	//return 24
	if sg.fs <= 14 {
		fh = sg.fs + int(font.Size)
	} else if sg.fs <= 20 {
		fh = sg.fs + 2*int(font.Size)
	} else {
		fh = sg.fs + 3*int(font.Size)
	}

	if fh == 0 {
		fh = 12
	}
	return
}

func (sg *OpenGLGraphics) FontMetrics(font chart.Font) (fw float32, fh int, mono bool) {
	//log.Panicf("Unimplemented: %s", whoami())
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

func (sg *OpenGLGraphics) TextLen(t string, font chart.Font) int {
	log.Panicf("Unimplemented: %s", whoami())
	// return chart.GenericTextLen(sg, t, font)
	return 0
}

var dashlength [][]int = [][]int{[]int{}, []int{4, 1}, []int{1, 1}, []int{4, 1, 1, 1, 1, 1}, []int{4, 4}, []int{1, 3}}

func (sg *OpenGLGraphics) Line(x0, y0, x1, y1 int, style chart.Style) {
	//log.Panicf("Unimplemented: %s", whoami())
	defer glh.OpenGLSentinel()()

	gl.LineWidth(float32(style.LineWidth))

	// TODO: line stipple?
	sc := chart.Color2RGBA(style.LineColor, uint8(style.Alpha*255))
	//log.Printf("color: %s %d %d %d %d", style.FillColor, sc.R, sc.G, sc.B, sc.A)

	gl.Color4ub(sc.R, sc.G, sc.B, sc.A)

	glh.With(glh.Attrib{gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT}, func() {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		glh.With(glh.Primitive{gl.LINES}, func() {
			gl.Vertex2i(x0, y0)
			gl.Vertex2i(x1, y1)
		})
	})
}

func (sg *OpenGLGraphics) Text(x, y int, t string, align string, rot int, f chart.Font) {
	if len(align) == 1 {
		align = "c" + align
	}

	_, fh, _ := sg.FontMetrics(f)
	tex := glh.MakeText(t, float64(fh))
	tex.Flipped = true
	defer tex.Destroy()

	switch align[0] {
	case 'b':
	case 'c':
		y += fh / 2
	case 't':
		y += fh
	default:
		log.Panicf("Unknown alignment: ", align)
	}

	switch align[1] {
	case 'l':
	case 'c':
		x -= tex.W / 2
	case 'r':
		x -= tex.W
	default:
		log.Panicf("Unknown alignment: ", align)
	}

	if f.Color != "" {
		sc := chart.Color2RGBA(f.Color, 0xff)
		gl.Color4ub(sc.R, sc.G, sc.B, sc.A)
	} else {
		// TODO: apply a default color?
	}
	glh.With(tex, func() {
		tex.Draw(x, y)
	})
}

func (sg *OpenGLGraphics) Symbol(x, y int, style chart.Style) {
	log.Panicf("Unimplemented: %s", whoami())
	// st := ""
	// filled := "fill:solid"
	// empty := "fill:none"
	if style.SymbolColor != "" {
		//st += "stroke:" + style.SymbolColor
		//filled = "fill:" + style.SymbolColor
		style.LineColor = style.SymbolColor
		style.Alpha = 1
	}
	f := style.SymbolSize
	if f == 0 {
		f = 1
	}
	// lw := 1
	// if style.LineWidth > 1 {
	// 	lw = style.LineWidth
	// }

	const n = 5         // default size
	a := int(n*f + 0.5) // standard
	// b := int(n/2*f + 0.5)     // smaller
	// c := int(1.155*n*f + 0.5) // triangel long sist
	// d := int(0.577*n*f + 0.5) // triangle short dist
	e := int(0.866*n*f + 0.5) // diagonal

	// sg.svg.Gstyle(fmt.Sprintf("%s; stroke-width: %d", st, lw))
	switch style.Symbol {
	case '*':
		sg.Line(x-e, y-e, x+e, y+e, style)
		sg.Line(x-e, y+e, x+e, y-e, style)
		fallthrough
	case '+':
		sg.Line(x-a, y, x+a, y, style)
		sg.Line(x, y-a, x, y+a, style)
	case 'X':
		sg.Line(x-e, y-e, x+e, y+e, style)
		sg.Line(x-e, y+e, x+e, y-e, style)
	case 'o':
		//sg.Circle(x, y, a, empty)
		panic("unimplemented")
	case '0':
		//sg.Circle(x, y, a, empty)
		panic("unimplemented")
		//sg.Circle(x, y, b, empty)
		panic("unimplemented")
	case '.':
		//sg.Circle(x, y, b, empty)
		panic("unimplemented")
	case '@':
		//sg.Circle(x, y, a, filled)
		panic("unimplemented")
	case '=':
		sg.Rect(x-e, y-e, 2*e, 2*e, style)
	case '#':
		sg.Rect(x-e, y-e, 2*e, 2*e, style)
	case 'A':
		//sg.Polygon([]int{x - a, x + a, x}, []int{y + d, y + d, y - c}, filled)
		panic("unimplemented")
	case '%':
		//sg.Polygon([]int{x - a, x + a, x}, []int{y + d, y + d, y - c}, empty)
		panic("unimplemented")
	case 'W':
		//sg.Polygon([]int{x - a, x + a, x}, []int{y - c, y - c, y + d}, filled)
		panic("unimplemented")
	case 'V':
		//sg.Polygon([]int{x - a, x + a, x}, []int{y - c, y - c, y + d}, empty)
		panic("unimplemented")
	case 'Z':
		//sg.Polygon([]int{x - e, x, x + e, x}, []int{y, y + e, y, y - e}, filled)
		panic("unimplemented")
	case '&':
		//sg.Polygon([]int{x - e, x, x + e, x}, []int{y, y + e, y, y - e}, empty)
		panic("unimplemented")
	default:
		panic("unimplemented")
		//sg.Text(x, y, "?", "text-anchor:middle; alignment-baseline:middle")
	}
	// sg.svg.Gend()

}

func (sg *OpenGLGraphics) Rect(x, y, w, h int, style chart.Style) {
	// log.Panicf("Unimplemented: %s", whoami())
	x, y, w, h = chart.SanitizeRect(x, y, w, h, style.LineWidth)
	defer glh.OpenGLSentinel()()

	//

	if style.FillColor != "" {
		fc := chart.Color2RGBA(style.FillColor, uint8(style.Alpha*255))

		//gl.Color4f(0.5, 0.25, 0.25, 1)
		//log.Printf("Fill color: %s %d %d %d %d", style.FillColor, fc.R, fc.G, fc.B, fc.A)

		glh.With(glh.Attrib{gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT}, func() {
			gl.Color4ub(fc.R, fc.G, fc.B, fc.A)

			gl.Enable(gl.BLEND)
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

			gl.Begin(gl.QUADS)
			glh.Squarei(x, y, w, h)
			gl.End()
		})
	}

	if style.LineWidth != 0 {
		gl.LineWidth(float32(style.LineWidth))
		//log.Print("Linewidth: ", float32(style.LineWidth))

		sc := chart.Color2RGBA(style.LineColor, uint8(style.Alpha*255))

		glh.With(glh.Attrib{gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT}, func() {
			gl.Color4ub(sc.R, sc.G, sc.B, sc.A)

			gl.Enable(gl.BLEND)
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

			gl.Begin(gl.LINE_LOOP)
			glh.Squarei(x, y, w, h)
			gl.End()
		})
	}

	// linecol := style.LineColor
	// if linecol != "" {
	// 	s = fmt.Sprintf("stroke:%s; ", linecol)
	// } else {
	// 	linecol = "#808080"
	// }
	// s += fmt.Sprintf("stroke-width: %d; ", style.LineWidth)
	// s += fmt.Sprintf("opacity: %.2f; ", 1-style.Alpha)
	// if style.FillColor != "" {
	// 	s += fmt.Sprintf("fill: %s; fill-opacity: %.2f", style.FillColor, 1-style.Alpha)
	// } else {
	// 	s += "fill-opacity: 0"
	// }
	// sg.svg.Rect(x, y, w, h, s)
	// GenericRect(sg, x, y, w, h, style) // TODO
}

func (sg *OpenGLGraphics) Path(x, y []int, style chart.Style) {
	log.Panicf("Unimplemented: %s", whoami())
	// n := len(x)
	// if len(y) < n {
	// 	n = len(y)
	// }
	// path := fmt.Sprintf("M %d,%d", x[0], y[0])
	// for i := 1; i < n; i++ {
	// 	path += fmt.Sprintf("L %d,%d", x[i], y[i])
	// }
	// st := linestyle(style)
	// sg.svg.Path(path, st)
}

func (sg *OpenGLGraphics) Wedge(x, y, ro, ri int, phi, psi float64, style chart.Style) {
	panic("No Wedge() for OpenGLGraphics.")
}

func (sg *OpenGLGraphics) XAxis(xr chart.Range, ys, yms int, options chart.PlotOptions) {
	//log.Panicf("Unimplemented: %s", whoami())
	chart.GenericXAxis(sg, xr, ys, yms, options)
	//log.Printf("X: %v %v %v %+v", xr, ys, yms, options)
}
func (sg *OpenGLGraphics) YAxis(yr chart.Range, xs, xms int, options chart.PlotOptions) {
	//log.Panicf("Unimplemented: %s", whoami())
	chart.GenericYAxis(sg, yr, xs, xms, options)
	//log.Printf("Y: %v %v %v %+v", yr, xs, xms, options)
}

func linestyle(style chart.Style) (s string) {
	log.Panicf("Unimplemented: %s", whoami())
	// lw := style.LineWidth
	// if style.LineColor != "" {
	// 	s = fmt.Sprintf("stroke:%s; ", style.LineColor)
	// }
	// s += fmt.Sprintf("stroke-width: %d; fill:none; ", lw)
	// s += fmt.Sprintf("opacity: %.2f; ", 1-style.Alpha)
	// if style.LineStyle != chart.SolidLine {
	// 	s += fmt.Sprintf("stroke-dasharray:")
	// 	for _, d := range dashlength[style.LineStyle] {
	// 		s += fmt.Sprintf(" %d", d*lw)
	// 	}
	// }
	return
}

func (sg *OpenGLGraphics) Scatter(points []chart.EPoint, plotstyle chart.PlotStyle, style chart.Style) {
	//log.Panicf("Unimplemented: %s", whoami())
	//chart.GenericScatter(sg, points, plotstyle, style)

	// TODO: Implement error bars/symbols

	var vertices glh.ColorVertices

	sc := chart.Color2RGBA(style.LineColor, uint8(style.Alpha*255))
	c := glh.Color{sc.R, sc.G, sc.B, sc.A}

	for _, p := range points {
		vertices.Add(glh.ColorVertex{c, glh.Vertex{float32(p.X), float32(p.Y)}})
	}

	gl.LineWidth(float32(style.LineWidth))

	glh.With(glh.Attrib{gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT}, func() {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		vertices.Draw(gl.LINE_STRIP)
	})

	/***********************************************
	// First pass: Error bars
	ebs := style
	ebs.LineColor, ebs.LineWidth, ebs.LineStyle = ebs.FillColor, 1, chart.SolidLine
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
			sg.Line(int(xl), int(p.Y), int(xh), int(p.Y), ebs)
		}
		if !math.IsNaN(p.DeltaY) {
			sg.Line(int(p.X), int(yl), int(p.X), int(yh), ebs)
		}
	}

	// Second pass: Line
	if (plotstyle&chart.PlotStyleLines) != 0 && len(points) > 0 {
		path := fmt.Sprintf("M %d,%d", int(points[0].X), int(points[0].Y))
		for i := 1; i < len(points); i++ {
			path += fmt.Sprintf("L %d,%d", int(points[i].X), int(points[i].Y))
		}
		st := linestyle(style)
		sg.svg.Path(path, st)
	}

	// Third pass: symbols
	if (plotstyle&chart.PlotStylePoints) != 0 && len(points) != 0 {
		for _, p := range points {
			sg.Symbol(int(p.X), int(p.Y), style)
		}
	}

	****************************************************/
}

func (sg *OpenGLGraphics) Boxes(boxes []chart.Box, width int, style chart.Style) {
	//log.Panicf("Unimplemented: %s", whoami())
	chart.GenericBoxes(sg, boxes, width, style)
}

func (sg *OpenGLGraphics) Key(x, y int, key chart.Key, options chart.PlotOptions) {
	//log.Panicf("Unimplemented: %s", whoami())
	chart.GenericKey(sg, x, y, key, options)
}

func (sg *OpenGLGraphics) Bars(bars []chart.Barinfo, style chart.Style) {
	//log.Panicf("Unimplemented: %s", whoami())
	chart.GenericBars(sg, bars, style)
}

func (sg *OpenGLGraphics) Rings(wedges []chart.Wedgeinfo, x, y, ro, ri int) {
	log.Panicf("Unimplemented: %s", whoami())
	// for _, w := range wedges {
	// 	var s string
	// 	linecol := w.Style.LineColor
	// 	if linecol != "" {
	// 		s = fmt.Sprintf("stroke:%s; ", linecol)
	// 	} else {
	// 		linecol = "#808080"
	// 	}
	// 	s += fmt.Sprintf("stroke-width: %d; ", w.Style.LineWidth)
	// 	s += fmt.Sprintf("opacity: %.2f; ", 1-w.Style.Alpha)
	// 	var sf string
	// 	if w.Style.FillColor != "" {
	// 		sf = fmt.Sprintf("fill: %s; fill-opacity: %.2f", w.Style.FillColor, 1-w.Style.Alpha)
	// 	} else {
	// 		sf = "fill-opacity: 0"
	// 	}

	// 	if math.Abs(w.Phi-w.Psi) >= 4*math.Pi {
	// 		sg.svg.Circle(x, y, ro, s+sf)
	// 		if ri > 0 {
	// 			sf = "fill: #ffffff; fill-opacity: 1"
	// 			sg.svg.Circle(x, y, ri, s+sf)
	// 		}
	// 		continue
	// 	}

	// 	var d string
	// 	p := 0.4 * float64(w.Style.LineWidth+w.Shift)
	// 	cphi, sphi := math.Cos(w.Phi), math.Sin(w.Phi)
	// 	cpsi, spsi := math.Cos(w.Psi), math.Sin(w.Psi)

	// 	if ri <= 0 {
	// 		// real wedge drawn as center -> outer radius -> arc -> closed to center
	// 		rf := float64(ro)
	// 		a := math.Sin((w.Psi - w.Phi) / 2)
	// 		dx, dy := p*math.Cos((w.Phi+w.Psi)/2)/a, p*math.Sin((w.Phi+w.Psi)/2)/a
	// 		d = fmt.Sprintf("M %d,%d ", x+int(dx+0.5), y+int(dy+0.5))

	// 		dx, dy = p*math.Cos(w.Phi+math.Pi/2), p*math.Sin(w.Phi+math.Pi/2)
	// 		d += fmt.Sprintf("L %d,%d ", int(rf*cphi+0.5+dx)+x, int(rf*sphi+0.5+dy)+y)

	// 		dx, dy = p*math.Cos(w.Psi-math.Pi/2), p*math.Sin(w.Psi-math.Pi/2)
	// 		d += fmt.Sprintf("A %d,%d 0 0 1 %d,%d ", ro, ro, int(rf*cpsi+0.5+dx)+x, int(rf*spsi+0.5+dy)+y)
	// 		d += fmt.Sprintf("z")
	// 	} else {
	// 		// ring drawn as inner radius -> outer radius -> outer arc -> inner radius -> inner arc
	// 		rof, rif := float64(ro), float64(ri)
	// 		dx, dy := p*math.Cos(w.Phi+math.Pi/2), p*math.Sin(w.Phi+math.Pi/2)
	// 		a, b := int(rif*cphi+0.5+dx)+x, int(rif*sphi+0.5+dy)+y
	// 		d = fmt.Sprintf("M %d,%d ", a, b)
	// 		d += fmt.Sprintf("L %d,%d ", int(rof*cphi+0.5+dx)+x, int(rof*sphi+0.5+dy)+y)

	// 		dx, dy = p*math.Cos(w.Psi-math.Pi/2), p*math.Sin(w.Psi-math.Pi/2)
	// 		d += fmt.Sprintf("A %d,%d 0 0 1 %d,%d ", ro, ro, int(rof*cpsi+0.5+dx)+x, int(rof*spsi+0.5+dy)+y)
	// 		d += fmt.Sprintf("L %d,%d ", int(rif*cpsi+0.5+dx)+x, int(rif*spsi+0.5+dy)+y)
	// 		d += fmt.Sprintf("A %d,%d 0 0 0 %d,%d ", ri, ri, a, b)
	// 		d += fmt.Sprintf("z")

	// 	}

	// 	sg.svg.Path(d, s+sf)

	// 	if w.Text != "" {
	// 		_, fh, _ := sg.FontMetrics(w.Font)
	// 		alpha := (w.Phi + w.Psi) / 2
	// 		var rt int
	// 		if ri > 0 {
	// 			rt = (ri + ro) / 2
	// 		} else {
	// 			rt = ro - 3*fh
	// 			if rt <= ro/2 {
	// 				rt = ro - 2*fh
	// 			}
	// 		}
	// 		tx, ty := int(float64(rt)*math.Cos(alpha)+0.5)+x, int(float64(rt)*math.Sin(alpha)+0.5)+y

	// 		sg.Text(tx, ty, w.Text, "cc", 0, w.Font)
	// 	}
	// }
}

var _ chart.Graphics = &OpenGLGraphics{}
