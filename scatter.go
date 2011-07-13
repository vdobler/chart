package chart

import (
	"math"
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
	Key            Key
	Data           []ScatterChartData
}

// Add any function f to this chart.
func (sc *ScatterChart) AddFunc(name string, f func(float64) float64) {
	s := Symbol[len(sc.Data)%len(Symbol)]
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, nil, f})
	ke := KeyEntry{Symbol:s, Text:name}
	sc.Key.Entries = append(sc.Key.Entries, ke)
}


// Add straight line through points (ax,ay) and (bx,by) tho chart.
func (sc *ScatterChart) AddLinear(name string, ax, ay, bx, by float64) {
	sc.AddFunc(name, func(x float64) float64 {
		return ay + (x-ax)*(by-ay)/(bx-ax)
	})
}


// Add points in data to chart.
func (sc *ScatterChart) AddData(name string, data []EPoint) {
	s := Symbol[len(sc.Data)%len(Symbol)]
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, data, nil})
	ke := KeyEntry{Symbol:s, Text:name}
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
	sc.AddData(name, edata)
}


// Make points from x and y and add to chart.
func (sc *ScatterChart) AddDataPair(name string, x, y []float64) {
	n := min(len(x), len(y))
	data := make([]EPoint, n)
	nan := math.NaN()
	for i := 0; i < n; i++ {
		data[i] = EPoint{X: x[i], Y: y[i], DeltaX: nan, DeltaY: nan}
	}
	sc.AddData(name, data)
}


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
	for s, data := range sc.Data {
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
				tb.Put(x, y, Symbol[s%len(Symbol)])
			}
		} else if data.Func != nil {
			// Functions. TODO(vodo) proper clipping
			symbol := Symbol[s%len(Symbol)]
			var lastsy, lastsx int
			var lastvalid bool
			for sx := leftm; sx < leftm+width; sx++ {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				sy := sc.YRange.Data2Screen(y)
				if y >= sc.YRange.Min && y <= sc.YRange.Max {
					if lastvalid {
						tb.Line(lastsx,lastsy, sx,sy, symbol)
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
