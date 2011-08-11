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
func (c *ScatterChart) AddFunc(name string, f func(float64) float64, plotstyle PlotStyle, style Style) {
	if plotstyle.undefined() {
		plotstyle = PlotStyleLines
	}
	if style.empty() {
		style = AutoStyle(len(c.Data), false)
	}

	scd := ScatterChartData{Name: name, PlotStyle: plotstyle, Style: style, Samples: nil, Func: f}
	c.Data = append(c.Data, scd)
	if name != "" {
		ke := KeyEntry{Text: name, PlotStyle: plotstyle, Style: style}
		c.Key.Entries = append(c.Key.Entries, ke)
	}
}


// AddData adds points in data to chart.
func (c *ScatterChart) AddData(name string, data []EPoint, plotstyle PlotStyle, style Style) {

	// Update styles if non given
	if plotstyle.undefined() {
		plotstyle = PlotStylePoints
	}
	if style.empty() {
		style = AutoStyle(len(c.Data), false)
	}

	// Init axis
	if len(c.Data) == 0 {
		c.XRange.init()
		c.YRange.init()
	}

	// Add data
	scd := ScatterChartData{Name: name, PlotStyle: plotstyle, Style: style, Samples: data, Func: nil}
	c.Data = append(c.Data, scd)

	// Autoscale
	for _, d := range data {
		xl, yl, xh, yh := d.BoundingBox()
		c.XRange.autoscale(xl)
		c.XRange.autoscale(xh)
		c.YRange.autoscale(yl)
		c.YRange.autoscale(yh)
	}

	// Add key/legend entry
	if name != "" {
		ke := KeyEntry{Style: style, PlotStyle: plotstyle, Text: name}
		c.Key.Entries = append(c.Key.Entries, ke)
	}
}


// Add points in data to chart.
func (c *ScatterChart) AddDataGeneric(name string, data []XYErrValue, plotstyle PlotStyle, style Style) {
	edata := make([]EPoint, len(data))
	for i, d := range data {
		x, y := d.XVal(), d.YVal()
		xl, xh := d.XErr()
		yl, yh := d.YErr()
		dx, dy := xh-xl, yh-yl
		xo, yo := xh-dx/2-x, yh-dy/2-y
		edata[i] = EPoint{X: x, Y: y, DeltaX: dx, DeltaY: dy, OffX: xo, OffY: yo}
	}
	c.AddData(name, edata, plotstyle, style)
}


// Make points from x and y and add to chart.
func (c *ScatterChart) AddDataPair(name string, x, y []float64, plotstyle PlotStyle, style Style) {
	n := min(len(x), len(y))
	data := make([]EPoint, n)
	nan := math.NaN()
	for i := 0; i < n; i++ {
		data[i] = EPoint{X: x[i], Y: y[i], DeltaX: nan, DeltaY: nan}
	}
	c.AddData(name, data, plotstyle, style)
}


// Plot outputs the scatter chart sc to g.
func (c *ScatterChart) Plot(g Graphics) {
	layout := Layout(g, c.Title, c.XRange.Label, c.YRange.Label,
		c.XRange.TicSetting.Hide, c.YRange.TicSetting.Hide, &c.Key)

	width, height := layout.Width, layout.Height
	topm, leftm := layout.Top, layout.Left
	numxtics, numytics := layout.NumXtics, layout.NumYtics

	fmt.Printf("\nSet up of X-Range (%d)\n", numxtics)
	c.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	fmt.Printf("\nSet up of Y-Range (%d)\n", numytics)
	c.YRange.Setup(numytics, numytics+2, height, topm, true)

	g.Begin()

	if c.Title != "" {
		g.Title(c.Title)
	}

	g.XAxis(c.XRange, topm+height, topm)
	g.YAxis(c.YRange, leftm, leftm+width)

	// Plot Data
	nan := math.NaN()
	xf, yf := c.XRange.Data2Screen, c.YRange.Data2Screen
	xmin, xmax := c.XRange.Min, c.XRange.Max
	ymin, ymax := c.YRange.Min, c.YRange.Max
	for _, data := range c.Data {
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
				x := c.XRange.Screen2Data(sx)
				y := data.Func(x)
				// TODO: half sample width if too f''(x) too big
				if y >= c.YRange.Min && y <= c.YRange.Max {
					sy := yf(y)
					p := EPoint{X: float64(sx), Y: float64(sy), DeltaX: nan, DeltaY: nan}
					points = append(points, p)
				} else {
					// TODO: buggy
					if y <= c.YRange.Min {
						sy := yf(c.YRange.Min)
						p := EPoint{X: float64(sx), Y: float64(sy), DeltaX: nan, DeltaY: nan}
						points = append(points, p)
					} else { // y > c.YRange.Max 
						sy := yf(c.YRange.Max)
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

	if !c.Key.Hide {
		g.Key(layout.KeyX, layout.KeyY, c.Key)
	}

	g.End()
}
