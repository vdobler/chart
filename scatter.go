package chart

import (
	"math"
)


// ScatterChart represents scatter charts and functions.
type ScatterChart struct {
	XRange, YRange Range  // x and y axis
	Title          string // Title of the chart
	Key            Key    // Key/Legend
	Data           []ScatterChartData
}

// ScatterChartData encapsulates a data set or function in a scatter chart.
// Not both Samples and Func may be non nil.
type ScatterChartData struct {
	Name    string
	Style   DataStyle
	Samples []EPoint
	Func    func(float64) float64
}

// AddFunc adds a function f to this chart.
func (sc *ScatterChart) AddFunc(name string, f func(float64) float64, style DataStyle) {
	if style.empty() {
		style = AutoStyle()
	}
	sc.Data = append(sc.Data, ScatterChartData{name, style, nil, f})
	ke := KeyEntry{Text: name, Style: style}
	sc.Key.Entries = append(sc.Key.Entries, ke)
}


// AddData adds points in data to chart.
func (sc *ScatterChart) AddData(name string, data []EPoint, style DataStyle) {
	if style.empty() {
		style = AutoStyle()
	}
	sc.Data = append(sc.Data, ScatterChartData{name, style, data, nil})
	ke := KeyEntry{Style: style, Text: name}
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


// Plot outputs the scatter chart sc to g.
func (sc *ScatterChart) Plot(g Graphics) {
	layout := Layout(g, sc.Title, sc.XRange.Label, sc.YRange.Label,
		sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)

	width, height := layout.Width, layout.Height
	topm, leftm := layout.Top, layout.Left
	numxtics, numytics := layout.NumXtics, layout.NumYtics

	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	sc.YRange.Setup(numytics, numytics+2, height, topm, true)

	g.Begin()

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
				if !math.IsNaN(d.DeltaX) {
					dx = float64(xf(d.DeltaX) - xf(0)) // TODO: abs?
					xo = float64(xf(d.OffX) - xf(0))
				}
				if !math.IsNaN(d.DeltaY) {
					dy = float64(yf(d.DeltaY) - yf(0)) // TODO: abs?
					yo = float64(yf(d.OffY) - yf(0))
				}
				// fmt.Printf("Point %d: %f\n", i, dx)
				p := EPoint{X: float64(x), Y: float64(y), DeltaX: dx, DeltaY: dy, OffX: xo, OffY: yo}
				points = append(points, p)
			}
			g.Scatter(points, style)
		} else if data.Func != nil {
			// Functions. TODO(vodo) proper clipping
			step := 8
			if width/step < 20 {
				step = 4
			}
			if width/step < 20 {
				step = 2
			}
			if width/step < 10 {
				step = 1
			}
			pcap := max(4, width/step)
			points := make([]EPoint, 0, pcap)

			for sx := leftm; sx < leftm+width; sx += step {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				// TODO: half sample width if too f''(x) too big
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

	if !sc.Key.Hide {
		g.Key(layout.KeyX, layout.KeyY, sc.Key)
	}

	g.End()
}
