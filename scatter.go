package chart

import (
	"math"
	"github.com/ajstarks/svgo"
	"io"
	"fmt"
	//	"os"
	//	"strings"
)


type ScatterChartData struct {
	Name    string
	Style   DataStyle
	Samples []EPoint
	Func    func(float64) float64
}


type ScatterChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Data           []ScatterChartData
}

// Add any function f to this chart.
func (sc *ScatterChart) AddFunc(name string, f func(float64) float64, style DataStyle) {
	if style.empty() {
		style = AutoStyle()
	}
	sc.Data = append(sc.Data, ScatterChartData{name, style, nil, f})
	ke := KeyEntry{Symbol: style.Symbol, Text: name}
	sc.Key.Entries = append(sc.Key.Entries, ke)
}


// Add points in data to chart.
func (sc *ScatterChart) AddData(name string, data []EPoint, style DataStyle) {
	if style.empty() {
		style = AutoStyle()
	}
	sc.Data = append(sc.Data, ScatterChartData{name, style, data, nil})
	ke := KeyEntry{Symbol: style.Symbol, Text: name}
	sc.Key.Entries = append(sc.Key.Entries, ke)
	if sc.XRange.DataMin == 0 && sc.XRange.DataMax == 0 && sc.YRange.DataMin == 0 && sc.YRange.DataMax == 0 {
		sc.XRange.DataMin = data[0].X
		sc.XRange.DataMax = data[0].X
		sc.YRange.DataMin = data[0].Y
		sc.YRange.DataMax = data[0].Y
	}
	for _, d := range data {
		xl, yl, xh, yh := d.boundingBox()
		if xl < sc.XRange.DataMin {
			sc.XRange.DataMin = xl
		} else if xh > sc.XRange.DataMax {
			sc.XRange.DataMax = xh
		}
		if yl < sc.YRange.DataMin {
			sc.YRange.DataMin = yl
		} else if yh > sc.YRange.DataMax {
			sc.YRange.DataMax = yh
		}
	}
	sc.XRange.Min = sc.XRange.DataMin
	sc.XRange.Max = sc.XRange.DataMax
	sc.YRange.Min = sc.YRange.DataMin
	sc.YRange.Max = sc.YRange.DataMax
}

// Add points in data to chart.
func (sc *ScatterChart) AddDataGeneric(name string, data []XYErrValue, style DataStyle) {
	edata := make([]EPoint, len(data))
	for i, d := range data {
		x, y := d.XVal(), d.YVal()
		xl, xh := d.XErr()
		yl, yh := d.YErr()
		dx, dy := xh-xl, yh-yl
		xo, yo := xh-dx/2-x, yh-dy/2-y
		edata[i] = EPoint{X: x, Y: y, DeltaX: dx, DeltaY: dy, OffX: xo, OffY: yo}
	}
	sc.AddData(name, edata, style)
}


// Make points from x and y and add to chart.
func (sc *ScatterChart) AddDataPair(name string, x, y []float64, style DataStyle) {
	n := min(len(x), len(y))
	data := make([]EPoint, n)
	nan := math.NaN()
	for i := 0; i < n; i++ {
		data[i] = EPoint{X: x[i], Y: y[i], DeltaX: nan, DeltaY: nan}
	}
	sc.AddData(name, data, style)
}


// PlotTxt renders the chat as ascii-art.
func (sc *ScatterChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key, 1, 1)

	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	sc.YRange.Setup(numytics, numytics+1, height, topm, true)

	tb := NewTextBuf(w, h)
	if sc.Title != "" {
		tb.Text(width/2+leftm, 0, sc.Title, 0)
	}

	TxtXRange(sc.XRange, tb, topm+height, topm, sc.Xlabel, 2)
	TxtYRange(sc.YRange, tb, leftm, leftm+width, sc.Ylabel, 2)

	// Plot Data
	nan := math.NaN()
	for _, data := range sc.Data {
		symbol := data.Style.Symbol
		if data.Samples != nil {
			// Samples
			for _, d := range data.Samples {
				x := sc.XRange.Data2Screen(d.X)
				y := sc.YRange.Data2Screen(d.Y)
				// TODO: clip
				if d.DeltaX != nan {
					xl, _, xh, _ := d.boundingBox()
					xe := sc.XRange.Data2Screen(xh)
					for xa := sc.XRange.Data2Screen(xl); xa <= xe; xa++ {
						tb.Put(xa, y, '-')
					}

				}
				if d.DeltaY != nan {
					_, yl, _, yh := d.boundingBox()
					ye := sc.YRange.Data2Screen(yh)
					for ya := sc.YRange.Data2Screen(yl); ya >= ye; ya-- {
						tb.Put(x, ya, '|')
					}

				}
				tb.Put(x, y, symbol)
			}
		} else if data.Func != nil {
			// Functions. TODO(vodo) proper clipping
			var lastsy, lastsx int
			var lastvalid bool
			for sx := leftm; sx < leftm+width; sx++ {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				sy := sc.YRange.Data2Screen(y)
				if y >= sc.YRange.Min && y <= sc.YRange.Max {
					if lastvalid {
						tb.Line(lastsx, lastsy, sx, sy, symbol)
					} else {
						lastvalid = true
					}
					lastsx, lastsy = sx, sy
				} else {
					lastvalid = false
				}
			}
		}
	}

	if kb != nil {
		tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	return tb.String()
}

// PlotSvg renders the chart as SVG
func (sc *ScatterChart) PlotSvg(w, h int, writer io.Writer) *svg.SVG {
	chart := svg.New(writer)
	chart.Start(w, h)
	chart.Title("SVG: " + sc.Title)
	var fontheight int
	switch {
	case w*h < 100*80:
		fontheight = 8
	case w*h < 200*120:
		fontheight = 9
	case w*h < 300*200:
		fontheight = 10
	case w*h < 400*300:
		fontheight = 11
	case w*h < 600*480:
		fontheight = 12
	case w*h < 800*600:
		fontheight = 13
	case w*h < 1024*800:
		fontheight = 14
	case w*h < 1400*1024:
		fontheight = 16
	default:
		fontheight = 18
	}
	fontwidth := (2 * fontheight) / 3
	chart.Gstyle(fmt.Sprintf("font-family: %s; font-size: %d", "Verdana", fontheight))
	chart.Rect(0, 0, w, h, "fill:white")

	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key, fontwidth, fontheight)

	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	sc.YRange.Setup(numytics, numytics+1, height, topm, true)

	if sc.Title != "" {
		col := "#000000"
		if pc, ok := Palette["title"]; ok {
			col = pc
		}
		style := fmt.Sprintf("text-anchor: middle; stroke: %s; font-size: %d", col, (12*fontheight)/10)
		chart.Text(width/2+leftm, 4*fontheight/3, sc.Title, style)
	}

	SvgXRange(sc.XRange, chart, topm+height, topm, sc.Xlabel, 2, fontwidth, fontheight)
	SvgYRange(sc.YRange, chart, leftm, leftm+width, sc.Ylabel, 2, fontwidth, fontheight)

	// Plot Data
	nan := math.NaN()
	for _, data := range sc.Data {
		style := data.Style
		if data.Samples != nil {
			// Samples
			for _, d := range data.Samples {
				x := sc.XRange.Data2Screen(d.X)
				y := sc.YRange.Data2Screen(d.Y)
				// TODO: clip
				if d.DeltaX != nan {
					xl, _, xh, _ := d.boundingBox()
					xa := sc.XRange.Data2Screen(xl)
					xe := sc.XRange.Data2Screen(xh)
					chart.Line(xa, y, xe, y, "stroke:gray; stroke-width:1")
				}
				if d.DeltaY != nan {
					_, yl, _, yh := d.boundingBox()
					ya := sc.YRange.Data2Screen(yl)
					ye := sc.YRange.Data2Screen(yh)
					chart.Line(x, ya, x, ye, "stroke:gray; stroke-width:1")
				}
				SvgSymbol(chart, x, y, style)
			}
		} else if data.Func != nil {
			// Functions. TODO(vodo) proper clipping
			// symbol := Symbol[s%len(Symbol)]
			var lastsy, lastsx int
			var lastvalid bool
			lineCol := style.LineColor
			if lineCol == "" {
				lineCol = style.SymbolColor
			}
			lineWidth := style.LineWidth
			if lineWidth == 0 {
				lineWidth = 2
			}
			st := fmt.Sprintf("stroke:%s; stroke-width:%d", lineCol, lineWidth)
			chart.Gstyle(st)
			for sx := leftm; sx < leftm+width; sx++ {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				sy := sc.YRange.Data2Screen(y)
				if y >= sc.YRange.Min && y <= sc.YRange.Max {
					if lastvalid {
						chart.Line(lastsx, lastsy, sx, sy)
					} else {
						lastvalid = true
					}
					lastsx, lastsy = sx, sy
				} else {
					lastvalid = false
				}
			}
			chart.Gend()
		}
	}

	if kb != nil {
		//	tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	chart.Gend()
	chart.End()
	return chart
}


func SvgSymbol(svg *svg.SVG, x, y int, style DataStyle) {
	s := style.Symbol
	col := style.SymbolColor
	f := style.SymbolSize
	if f == 0 {
		f = 1
	}
	const n = 5               // default size
	a := int(n*f + 0.5)       // standard
	b := int(n/2*f + 0.5)     // smaller
	c := int(1.155*n*f + 0.5) // triangel long sist
	d := int(0.577*n*f + 0.5) // triangle short dist
	e := int(0.866*n*f + 0.5) // diagonal

	svg.Gstyle("stroke:" + col + "; stroke-width: 1")
	switch s {
	case '*':
		svg.Line(x-e, y-e, x+e, y+e)
		svg.Line(x-e, y+e, x+e, y-e)
		fallthrough
	case '+':
		svg.Line(x-a, y, x+a, y)
		svg.Line(x, y-a, x, y+a)
	case 'X':
		svg.Line(x-e, y-e, x+e, y+e)
		svg.Line(x-e, y+e, x+e, y-e)
	case 'o':
		svg.Circle(x, y, a, "fill:none")
	case '0':
		svg.Circle(x, y, a, "fill:none")
		svg.Circle(x, y, b, "fill:none")
	case '.':
		svg.Circle(x, y, b, "fill:none")
	case '@':
		svg.Circle(x, y, a, "fill:"+col)
	case '=':
		svg.Rect(x-e, y-e, 2*e, 2*e, "fill:none")
	case '#':
		svg.Rect(x-e, y-e, 2*e, 2*e, "fill:"+col)
	case 'A':
		svg.Polygon([]int{x - a, x + a, x}, []int{y + d, y + d, y - c}, "fill:"+col)
	case '%':
		svg.Polygon([]int{x - a, x + a, x}, []int{y + d, y + d, y - c}, "fill:none")
	case 'W':
		svg.Polygon([]int{x - a, x + a, x}, []int{y - c, y - c, y + d}, "fill:"+col)
	case 'V':
		svg.Polygon([]int{x - a, x + a, x}, []int{y - c, y - c, y + d}, "fill:none")
	case 'Z':
		svg.Polygon([]int{x - e, x, x + e, x}, []int{y, y + e, y, y - e}, "fill:"+col)
	case '&':
		svg.Polygon([]int{x - e, x, x + e, x}, []int{y, y + e, y, y - e}, "fill:none")
	default:
		svg.Text(x, y, "?", "text-anchor:middle; alignment-baseline:middle")
	}
	svg.Gend()
}


// Print xrange to tb at vertical position y.
// Axis, tics, tic labels, axis label and range limits are drawn.
// mirror: 0: no other axis, 1: axis without tics, 2: axis with tics,
func SvgXRange(xrange Range, chart *svg.SVG, y, y1 int, label string, mirror int, fontwidth, fontheight int) {
	var ticLen int = 0
	if !xrange.TicSetting.Hide {
		ticLen = 5
	}

	xa, xe := xrange.Data2Screen(xrange.Min), xrange.Data2Screen(xrange.Max)
	chart.Line(xa, y, xe, y, "stroke:black; stroke-width:2")
	if mirror >= 1 {
		chart.Line(xa, y1, xe, y1, "stroke:black; stroke-width:2")
	}
	if xrange.ShowZero && xrange.Min < 0 && xrange.Max > 0 {
		z := xrange.Data2Screen(0)
		chart.Line(z, y, z, y1, "stroke:gray; stroke-width:1")
	}

	if label != "" {
		yy := y + fontheight
		if !xrange.TicSetting.Hide {
			yy += 3*fontheight/2 + ticLen
		}
		chart.Text((xa+xe)/2, yy, label, "text-anchor:middle")
	}

	for ticcnt, tic := range xrange.Tics {
		x := xrange.Data2Screen(tic.Pos)
		lx := xrange.Data2Screen(tic.LabelPos)
		if ticcnt > 0 && ticcnt < len(xrange.Tics)-1 && xrange.TicSetting.Grid == 1 {
			chart.Line(x, y-1, x, y1+1, "stroke: #808080; stroke-width:1")
		}
		chart.Line(x, y-ticLen, x, y+ticLen, "stroke:black; stroke-width:2")
		if mirror >= 2 {
			chart.Line(x, y1-ticLen, x, y1+ticLen, "stroke:black; stroke-width:2")
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

func SvgYRange(yrange Range, chart *svg.SVG, x, x1 int, label string, mirror int, fontwidth, fontheight int) {
	var ticLen int = 0
	if !yrange.TicSetting.Hide {
		ticLen = 5
	}

	ya, ye := yrange.Data2Screen(yrange.Min), yrange.Data2Screen(yrange.Max)
	chart.Line(x, ya, x, ye, "stroke:black; stroke-width:2")
	if mirror >= 1 {
		chart.Line(x1, ya, x1, ye, "stroke:black; stroke-width:2")
	}
	if yrange.ShowZero && yrange.Min < 0 && yrange.Max > 0 {
		z := yrange.Data2Screen(0)
		chart.Line(x, z, x1, z, "stroke:gray; stroke-width:1")
	}

	if label != "" {
		x := 2 * fontheight
		y := (ya+ye)/2 + fontheight/2
		trans := fmt.Sprintf("transform=\"rotate(-90 %d %d)\"", 16, (ya+ye)/2)
		chart.Text(x, y, label, trans, "text-anchor:middle;")
	}

	for _, tic := range yrange.Tics {
		y := yrange.Data2Screen(tic.Pos)
		ly := yrange.Data2Screen(tic.LabelPos) + 4*fontheight/10
		chart.Line(x-ticLen, y, x+ticLen, y, "stroke:black; stroke-width:2")
		if mirror >= 2 {
			chart.Line(x1-ticLen, y, x1+ticLen, y, "stroke:black; stroke-width:2")
		}
		if yrange.Time {
			if tic.Align == 0 { // centered tic
				chart.Line(x-2*ticLen, y, x-ticLen, y, "stroke:black; stroke-width:2")
			}
			chart.Text(x-ticLen-fontwidth, ly, tic.Label,
				"text-anchor:end; alignment-baseline:middle")
		} else {

			chart.Text(x-ticLen-fontwidth, ly, tic.Label, "text-anchor:end;")
		}
	}
}


func (sc *ScatterChart) Plot(g Graphics) {
	fontwidth, fontheight := g.FontMetrics()
	w, h := g.Dimensions()

	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key, fontwidth, fontheight)


	g.Begin()

	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	sc.YRange.Setup(numytics, numytics+1, height, topm, true)

	if sc.Title != "" {
		g.Title(sc.Title)
	}

	g.XAxis(sc.XRange, topm+height, topm)
	g.YAxis(sc.YRange, leftm, leftm+width)

	// Plot Data
	nan := math.NaN()
	xf, yf := sc.XRange.Data2Screen, sc.YRange.Data2Screen
	for _, data := range sc.Data {
		style := data.Style
		if data.Samples != nil {
			// Samples
			points := make([]EPoint, 0, len(data.Samples))
			for _, d := range data.Samples {
				x := xf(d.X)
				y := yf(d.Y)
				dx, dy := nan, nan
				var xo, yo float64
				// TODO: clip
				if d.DeltaX != nan {
					dx = float64(xf(d.DeltaX) - xf(0)) // TODO: abs?
					xo = float64(xf(d.OffX) - xf(0))
				}
				if d.DeltaY != nan {
					dy = float64(yf(d.DeltaY) - yf(0)) // TODO: abs?
					yo = float64(yf(d.OffY) - yf(0))
				}
				p := EPoint{X: float64(x), Y: float64(y), DeltaX: dx, DeltaY: dy, OffX: xo, OffY: yo}
				points = append(points, p)
			}
			g.Scatter(points, style)
		} else if data.Func != nil {
			// Functions. TODO(vodo) proper clipping
			points := make([]EPoint, 0, width/10)

			for sx := leftm; sx < leftm+width; sx+=10 {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				if y >= sc.YRange.Min && y <= sc.YRange.Max {
					sy := yf(y)
					p := EPoint{X: float64(sx), Y: float64(sy), DeltaX: nan, DeltaY: nan}
					points = append(points, p)
				} else {
					// TODO: buggy
					if y <= sc.YRange.Min {
						sy := yf(sc.YRange.Min)
						p := EPoint{X: float64(sx), Y: float64(sy), DeltaX: nan, DeltaY: nan}
						points = append(points, p)
					} else { // y > sc.YRange.Max 
						sy := yf(sc.YRange.Max)
						p := EPoint{X: float64(sx), Y: float64(sy), DeltaX: nan, DeltaY: nan}
						points = append(points, p)
					}
					g.Scatter(points, style)
					points = make([]EPoint, 0, width/10)
				}
			}
			g.Scatter(points, style)
		}
	}

	if kb != nil {
		//	tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	g.End()
}