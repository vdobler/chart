package chart

import (
	"fmt"
	"math"
)


// ScatterChart represents scatter charts, line charts and function plots.
type ScatterChart struct {
	XRange, YRange Range              // X and Y axis
	Title          string             // Title of the chart
	Key            Key                // Key/Legend
	Data           []ScatterChartData // The actual data (filled with Add...-methods)
}

// ScatterChartData encapsulates a data set or function in a scatter chart.
// Not both Samples and Func may be non nil at the same time.
type ScatterChartData struct {
	Name      string                // The name of this data set. TODO: unused?
	PlotStyle PlotStyle             // Points, Lines+Points or Lines only
	Style     Style                 // Color, sizes, pointtype, linestyle, ...
	Samples   []EPoint              // The actual points for scatter/lines charts
	Func      func(float64) float64 // The function to draw.
}

// AddFunc adds a function f to this chart.
func (sc *ScatterChart) AddFunc(name string, f func(float64) float64, plotstyle PlotStyle, style Style) {
	if plotstyle.undefined() {
		plotstyle = PlotStyleLines
	}
	if style.empty() {
		style = AutoStyle(len(sc.Data), false)
	}

	scd := ScatterChartData{Name: name, PlotStyle: plotstyle, Style: style, Samples: nil, Func: f}
	sc.Data = append(sc.Data, scd)
	if name != "" {
		ke := KeyEntry{Text: name, PlotStyle: plotstyle, Style: style}
		sc.Key.Entries = append(sc.Key.Entries, ke)
	}
}


// AddData adds points in data to chart.
func (sc *ScatterChart) AddData(name string, data []EPoint, plotstyle PlotStyle, style Style) {

	// Update styles if non given
	if plotstyle.undefined() {
		plotstyle = PlotStylePoints
	}
	if style.empty() {
		style = AutoStyle(len(sc.Data), false)
	}

	// Init axis
	if len(sc.Data) == 0 {
		sc.XRange.init()
		sc.YRange.init()
	}

	// Add data
	scd := ScatterChartData{Name: name, PlotStyle: plotstyle, Style: style, Samples: data, Func: nil}
	sc.Data = append(sc.Data, scd)

	// Add key/legend entry
	if name != "" {
		ke := KeyEntry{Style: style, PlotStyle: plotstyle, Text: name}
		sc.Key.Entries = append(sc.Key.Entries, ke)
	}

	// Autoscale
	for _, d := range data {
		xl, yl, xh, yh := d.BoundingBox()
		sc.XRange.autoscale(xl)
		sc.XRange.autoscale(xh)
		sc.YRange.autoscale(yl)
		sc.YRange.autoscale(yh)
	}

}


func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Prepare the range r for use, especially set up all values needed for autoscale() to work properly
func (r *Range) init() {
	// All the min stuff
	if r.MinMode.Fixed {
		// copy TValue to Value if set and time axis
		if r.Time && r.MinMode.TValue != nil {
			r.MinMode.Value = float64(r.MinMode.TValue.Seconds())
		}
		r.DataMin = r.MinMode.Value
	} else if r.MinMode.Constrained {
		// copy TLower/TUpper to Lower/Upper if set and time axis
		if r.Time && r.MinMode.TLower != nil {
			r.MinMode.Lower = float64(r.MinMode.TLower.Seconds())
		}
		if r.Time && r.MinMode.TUpper != nil {
			r.MinMode.Upper = float64(r.MinMode.TUpper.Seconds())
		}
		if r.MinMode.Lower == 0 && r.MinMode.Upper == 0 {
			// Constrained but un-initialized: Full autoscaling
			r.MinMode.Lower = -math.MaxFloat64
			r.MinMode.Upper = math.MaxFloat64
		}
		r.DataMin = r.MinMode.Upper
	} else {
		r.DataMin = math.MaxFloat64
	}

	// All the max stuff
	if r.MaxMode.Fixed {
		// copy TValue to Value if set and time axis
		if r.Time && r.MaxMode.TValue != nil {
			r.MaxMode.Value = float64(r.MaxMode.TValue.Seconds())
		}
		r.DataMax = r.MaxMode.Value
	} else if r.MaxMode.Constrained {
		// copy TLower/TUpper to Lower/Upper if set and time axis
		if r.Time && r.MaxMode.TLower != nil {
			r.MaxMode.Lower = float64(r.MaxMode.TLower.Seconds())
		}
		if r.Time && r.MaxMode.TUpper != nil {
			r.MaxMode.Upper = float64(r.MaxMode.TUpper.Seconds())
		}
		if r.MaxMode.Lower == 0 && r.MaxMode.Upper == 0 {
			// Constrained but un-initialized: Full autoscaling
			r.MaxMode.Lower = -math.MaxFloat64
			r.MaxMode.Upper = math.MaxFloat64
		}
		r.DataMax = r.MaxMode.Upper
	} else {
		r.DataMax = -math.MaxFloat64
	}

	fmt.Printf("At end of init: DataMin / DataMax  =   %g / %g\n", r.DataMin, r.DataMax)
}


// Update DataMin and DataMax according to the RangeModes.
func (r *Range) autoscale(x float64) {

	if x < r.DataMin && !r.MinMode.Fixed {
		if !r.MinMode.Constrained {
			// full autoscaling
			r.DataMin = x
		} else {
			r.DataMin = fmin(fmax(x, r.MinMode.Lower), r.DataMin)
		}
	}

	if x > r.DataMax && !r.MaxMode.Fixed {
		if !r.MaxMode.Constrained {
			// full autoscaling
			r.DataMax = x
		} else {
			r.DataMax = fmax(fmin(x, r.MaxMode.Upper), r.DataMax)
		}
	}
}

/*
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
*/

// Make points from x and y and add to chart.
func (sc *ScatterChart) AddDataPair(name string, x, y []float64, plotstyle PlotStyle, style Style) {
	n := min(len(x), len(y))
	data := make([]EPoint, n)
	nan := math.NaN()
	for i := 0; i < n; i++ {
		data[i] = EPoint{X: x[i], Y: y[i], DeltaX: nan, DeltaY: nan}
	}
	sc.AddData(name, data, plotstyle, style)
}


// Plot outputs the scatter chart sc to g.
func (sc *ScatterChart) Plot(g Graphics) {
	layout := Layout(g, sc.Title, sc.XRange.Label, sc.YRange.Label,
		sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)

	width, height := layout.Width, layout.Height
	topm, leftm := layout.Top, layout.Left
	numxtics, numytics := layout.NumXtics, layout.NumYtics

	fmt.Printf("\nSet up of X-Range (%d)\n", numxtics)
	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	fmt.Printf("\nSet up of Y-Range (%d)\n", numytics)
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
	xmin, xmax := sc.XRange.Min, sc.XRange.Max
	ymin, ymax := sc.YRange.Min, sc.YRange.Max
	for _, data := range sc.Data {
		style := data.Style
		if data.Samples != nil {
			// Samples
			points := make([]EPoint, 0, len(data.Samples))
			for _, d := range data.Samples {
				if d.X < xmin || d.X > xmax || d.Y < ymin || d.Y > ymax {
					continue
				}
				xl, yl, xh, yh := d.BoundingBox()
				if xl < xmin { // happens only if d.Delta!=0,NaN
					a := xmin - xl
					d.DeltaX -= a
					d.OffX += a / 2
				}
				if xh > xmax {
					a := xh - xmax
					d.DeltaX -= a
					d.OffX -= a / 2
				}
				if yl < ymin { // happens only if d.Delta!=0,NaN
					a := ymin - yl
					d.DeltaY -= a
					d.OffY += a / 2
				}
				if yh > ymax {
					a := yh - ymax
					d.DeltaY -= a
					d.OffY -= a / 2
				}

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
			g.Scatter(points, data.PlotStyle, style)
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
					g.Scatter(points, data.PlotStyle, style)
					points = make([]EPoint, 0, width/10)
				}
			}
			g.Scatter(points, data.PlotStyle, style)
		}
	}

	if !sc.Key.Hide {
		g.Key(layout.KeyX, layout.KeyY, sc.Key)
	}

	g.End()
}
