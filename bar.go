package chart

import (
	"fmt"
	"math"
	//	"os"
	//	"strings"
)


type BarChartData struct {
	Name    string
	Style   DataStyle
	Samples []Point
}


// Simple Bar Chart
//  
// x values must be sorted.
// 
// In stacked mode all x values of all data sets must be identical. Not even
// missing values are allowed.
// 
// If BarWidth is zero the BarWidth is the smallest distance between two
// x values multiplied ny BarWidthFac (<1). 
// Data is drwan first to last, last overwriting previous, maybe
// at the same x position.  Use CategoricalBarChart if your x axis is
// categorical.
//

type BarChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Horizontal     bool    // Display is horizontal bars
	Stacked        bool    // Display different data sets ontop of each other
	ShowVal        bool    // Display values 
	SameBarWidth   bool    // all data sets use smalest bar width
	BarWidthFac    float64 // if nonzero: use this as width for all bars
	Data           []BarChartData
}

func (c *BarChart) AddData(name string, data []Point) {
	s := Symbol[len(c.Data)%len(Symbol)]
	c.Data = append(c.Data, BarChartData{name, DataStyle{Symbol: s}, data})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{s, name})

	if len(c.Data) == 1 { // first data set 
		c.XRange.DataMin = data[0].X
		c.XRange.DataMax = data[0].X
		c.YRange.DataMin = data[0].Y
		c.YRange.DataMax = data[0].Y
	}
	for _, d := range data {
		if d.X < c.XRange.DataMin {
			c.XRange.DataMin = d.X
		} else if d.X > c.XRange.DataMax {
			c.XRange.DataMax = d.X
		}
		if d.Y < c.YRange.DataMin {
			c.YRange.DataMin = d.Y
		} else if d.Y > c.YRange.DataMax {
			c.YRange.DataMax = d.Y
		}
	}
	c.XRange.Min = c.XRange.DataMin
	c.XRange.Max = c.XRange.DataMax
}

func (c *BarChart) AddDataPair(name string, x, y []float64) {
	n := min(len(x), len(y))
	data := make([]Point, n)
	for i := 0; i < n; i++ {
		data[i] = Point{X: x[i], Y: y[i]}
	}
	c.AddData(name, data)
}

func (c *BarChart) barWidth(sample int) float64 {
	// find bar width
	barWidth := c.XRange.Max - c.XRange.Min // large enough
	for i := 1; i < len(c.Data[sample].Samples); i++ {
		diff := math.Fabs(c.Data[sample].Samples[i].X - c.Data[sample].Samples[i-1].X)
		if diff < barWidth {
			barWidth = diff
		}
	}
	if c.BarWidthFac != 0 {
		barWidth *= math.Fabs(c.BarWidthFac)
	}
	fmt.Printf("BarWidth for sample %d: %f\n", sample, barWidth)
	return barWidth
}

func (c *BarChart) extremBarWidth() (smallest, widest float64) {
	w0 := c.barWidth(0)
	widest, smallest = w0, w0
	for s := 1; s < len(c.Data); s++ {
		b := c.barWidth(s)
		if b > widest {
			widest = b
		} else if b < smallest {
			smallest = b
		}
	}
	return
}


func (c *BarChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, c.Title, c.Xlabel, c.Ylabel, c.XRange.TicSetting.Hide, c.YRange.TicSetting.Hide, &c.Key)

	// Outside bound ranges for histograms are nicer
	leftm, width = leftm+2, width-2
	topm, height = topm, height-1

	// find bar width
	lbw, ubw := c.extremBarWidth()
	var barWidth float64
	if c.SameBarWidth {
		barWidth = lbw
	} else {
		barWidth = ubw
	}

	// set up range and extend if bar would not fit
	c.XRange.Setup(numxtics, numxtics+1, width, leftm, false)
	c.YRange.Setup(numytics, numytics+2, height, topm, true)

	if c.XRange.DataMin-barWidth/2 < c.XRange.Min {
		c.XRange.DataMin -= barWidth / 2
	}
	if c.XRange.DataMax+barWidth > c.XRange.Max {
		c.XRange.DataMax += barWidth / 2
	}
	c.XRange.Setup(numxtics, numxtics+1, width, leftm, false)

	tb := NewTextBuf(w, h)
	if c.Title != "" {
		tb.Text(width/2+leftm, 0, c.Title, 0)
	}
	TxtXRange(c.XRange, tb, topm+height+1, 0, c.Xlabel, 0)
	TxtYRange(c.YRange, tb, leftm-2, 0, c.Ylabel, 0)

	xf := c.XRange.Data2Screen
	yf := c.YRange.Data2Screen
	sy0 := yf(c.YRange.Min)

	barWidth = lbw
	for i, data := range c.Data {
		if !c.SameBarWidth {
			barWidth = c.barWidth(i)
		}
		sbw := max(1, xf(2*barWidth)-xf(barWidth)-1) // screen bar width
		symbol := data.Style.Symbol
		for _, point := range data.Samples {
			x, y := point.X, point.Y
			sx := xf(x-barWidth/2) + 1
			// sw := xf(x+barWidth/2) - sx
			sy := yf(y)
			sh := sy0 - sy
			tb.Block(sx, sy, sbw, sh, symbol)
		}
	}

	if kb != nil {
		tb.Paste(c.Key.X, c.Key.Y, kb)
	}

	return tb.String()
}
