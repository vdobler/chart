package chart

import (
	"math"
	"github.com/ajstarks/svgo"
	"io"
	// "fmt"
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
	Style          DataStyle
	Key            Key
	Data           []ScatterChartData
}

// Add any function f to this chart.
func (sc *ScatterChart) AddFunc(name string, f func(float64) float64) {
	s := Symbol[len(sc.Data)%len(Symbol)]
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, nil, f})
	ke := KeyEntry{Symbol: s, Text: name}
	sc.Key.Entries = append(sc.Key.Entries, ke)
}


// Add straight line through points (ax,ay) and (bx,by) tho chart.
func (sc *ScatterChart) AddLinear(name string, ax, ay, bx, by float64) {
	sc.AddFunc(name, func(x float64) float64 {
		return ay + (x-ax)*(by-ay)/(bx-ax)
	})
}


// Add points in data to chart.
func (sc *ScatterChart) AddData(name string, data []EPoint, style DataStyle) {
	sc.Style = sc.Style.NextMerge(style)
	sc.Data = append(sc.Data, ScatterChartData{name, sc.Style, data, nil})
	ke := KeyEntry{Symbol: Symbol[sc.Style.Symbol], Text: name}
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
func (sc *ScatterChart) AddDataGeneric(name string, data []XYErrValue) {
	edata := make([]EPoint, len(data))
	for i, d := range data {
		x, y := d.XVal(), d.YVal()
		xl, xh := d.XErr()
		yl, yh := d.YErr()
		dx, dy := xh-xl, yh-yl
		xo, yo := xh-dx/2-x, yh-dy/2-y
		edata[i] = EPoint{X: x, Y: y, DeltaX: dx, DeltaY: dy, OffX: xo, OffY: yo}
	}
	sc.AddData(name, edata, DataStyle{})
}


// Make points from x and y and add to chart.
func (sc *ScatterChart) AddDataPair(name string, x, y []float64) {
	n := min(len(x), len(y))
	data := make([]EPoint, n)
	nan := math.NaN()
	for i := 0; i < n; i++ {
		data[i] = EPoint{X: x[i], Y: y[i], DeltaX: nan, DeltaY: nan}
	}
	sc.AddData(name, data, DataStyle{})
}


// PlotTxt renders the chat as ascii-art.
func (sc *ScatterChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)

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
		symbol := Symbol[data.Style.Symbol]
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
	chart.Rect(0, 0, w, h, "fill:white")

	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)

	leftm += 100
	topm += 20
	width -= 150
	height -= 60

	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	sc.YRange.Setup(numytics, numytics+1, height, topm, true)

	if sc.Title != "" {
		chart.Text(width/2+leftm, 10, sc.Title, "")
	}

	SvgXRange(sc.XRange, chart, topm+height, topm, sc.Xlabel, 2)
	// SvgYRange(sc.YRange, chart, leftm, leftm+width, sc.Ylabel, 2)

	// Plot Data
	nan := math.NaN()
	for _, data := range sc.Data {
		if data.Samples != nil {
			// Samples
			symbol := Symbol[data.Style.Symbol]
			color := Palette[data.Style.SymbolColor]
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
				SvgSymbol(chart, x, y, symbol, color)
			}
		} else if data.Func != nil {
			// Functions. TODO(vodo) proper clipping
			// symbol := Symbol[s%len(Symbol)]
			var lastsy, lastsx int
			var lastvalid bool
			for sx := leftm; sx < leftm+width; sx++ {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				sy := sc.YRange.Data2Screen(y)
				if y >= sc.YRange.Min && y <= sc.YRange.Max {
					if lastvalid {
						chart.Line(lastsx, lastsy, sx, sy, "stroke:red; stroke-width:1") // TODO use style
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
		//	tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	chart.End()
	return chart
}


func SvgSymbol(svg *svg.SVG, x, y, s int, col string) {
	svg.Gstyle("stroke:" + col + "; stroke-width: 1")
	switch s {
	case '*':
		svg.Line(x-5, y-5, x+5, y+5)
		svg.Line(x-5, y+5, x+5, y-5)
		fallthrough
	case '+':
		svg.Line(x-5, y, x+5, y)
		svg.Line(x, y-5, x, y+5)
	case 'o':
		svg.Circle(x, y, 5)
	default:
		svg.Rect(x-5, y-5, 10, 10)
	}
	svg.Gend()
}

var fontheight = 12
var fontwidth = 8

// Print xrange to tb at vertical position y.
// Axis, tics, tic labels, axis label and range limits are drawn.
// mirror: 0: no other axis, 1: axis without tics, 2: axis with tics,
func SvgXRange(xrange Range, chart *svg.SVG, y, y1 int, label string, mirror int) {
	xa, xe := xrange.Data2Screen(xrange.Min), xrange.Data2Screen(xrange.Max)
	chart.Line(xa, y, xe, y, "stroke:black; strock-width:2")
	if mirror >= 1 {
		chart.Line(xa, y1, xe, y1, "stroke:black; strock-width:2")
	}
	if xrange.ShowZero && xrange.Min < 0 && xrange.Max > 0 {
		z := xrange.Data2Screen(0)
		chart.Line(z, y, z, y1, "stroke:gray; strocke-width:1")
	}

	if label != "" {
		yy := y + fontheight
		if !xrange.TicSetting.Hide {
			yy += fontheight
		}
		chart.Text((xa+xe)/2, yy, label, "text-anchor:middle")
	}

	for _, tic := range xrange.Tics {
		x := xrange.Data2Screen(tic.Pos)
		lx := xrange.Data2Screen(tic.LabelPos)
		if xrange.Time {
			chart.Line(x, y-5, x, y+15, "stroke:black; strocke-width:2")
			if mirror >= 2 {
				chart.Line(x, y1-5, x, y1+5, "stroke:black; strocke-width:2")
			}
			if tic.Align == -1 {
				chart.Text(lx, y+fontheight, tic.Label, "text-anchor:left")
			} else {
				chart.Text(lx, y+fontheight, tic.Label, "text-anchor:middle")
			}
		} else {
			chart.Line(x, y-5, x, y+5, "stroke:black; strocke-width:2")
			if mirror >= 2 {
				chart.Line(x, y1-5, x, y1+5, "stroke:black; strocke-width:2")
			}
			chart.Text(lx, y+fontheight, tic.Label, "text-anchor:middle")
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
}
