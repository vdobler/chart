package chart

import (
//"fmt"
//"math"
//	"os"
//	"strings"
)


type CategoryBarChartData struct {
	Name    string
	Style   DataStyle
	Samples map[string]float64
}


// Bar Chart with categorical x-axis
//  
type CategoryBarChart struct {
	Categories     []string // List of categories to display (ordered from left to right)
	YRange         Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Horizontal     bool // Display is horizontal bars
	Stacked        bool // Display different data sets ontop of each other
	ShowVal        bool // Display values 
	Data           []CategoryBarChartData
}

func (c *CategoryBarChart) AddData(name string, data map[string]float64) {
	s := Symbol[len(c.Data)%len(Symbol)]
	c.Data = append(c.Data, CategoryBarChartData{name, DataStyle{Symbol: s}, data})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{s, name})

	if len(c.Data) == 1 { // first data set
		for _, v := range data {
			c.YRange.DataMin = v
			c.YRange.DataMax = v
			break
		}
	}
	for _, v := range data {
		if v < c.YRange.DataMin {
			c.YRange.DataMin = v
		} else if v > c.YRange.DataMax {
			c.YRange.DataMax = v
		}
	}
	c.YRange.Min = c.YRange.DataMin
	c.YRange.Max = c.YRange.DataMax
}


func (c *CategoryBarChart) catIdx(cat string) (i int) {
	i = -1
	for n, c := range c.Categories {
		if c == cat {
			i = n
			return
		}
	}
	return
}

func (c *CategoryBarChart) PlotTxt(w, h int) string {
	n := len(c.Categories)
	xrange := Range{}
	xrange.DataMin, xrange.DataMax = 1, float64(n)
	xrange.Min, xrange.Max = 0.5, float64(n)+0.5
	xrange.MinMode = RangeMode{Fixed: true, Value: 0.5}
	xrange.MaxMode = RangeMode{Fixed: true, Value: float64(n) + 0.5}

	width, leftm, height, topm, kb, _, numytics := LayoutTxt(w, h, c.Title, c.Xlabel, c.Ylabel, xrange.TicSetting.Hide, c.YRange.TicSetting.Hide, &c.Key)

	// Outside bound ranges for histograms are nicer
	leftm, width = leftm+2, width-2
	topm, height = topm, height-1

	// find bar width
	var barWidth float64 = 0.3

	// set up range and extend if bar would not fit
	xrange.Setup(n, n, width, leftm, false)
	c.YRange.Setup(numytics, numytics+2, height, topm, true)
	xrange.Tics = make([]Tic, n)
	for i := 0; i < n; i++ {
		xrange.Tics[i].Pos = -1 // outside, no tic
		xrange.Tics[i].LabelPos = float64(i) + 1
		xrange.Tics[i].Label = c.Categories[i]
		xrange.Tics[i].Align = 0 // center
	}

	tb := NewTextBuf(w, h)
	if c.Title != "" {
		tb.Text(width/2+leftm, 0, c.Title, 0)
	}
	TxtXRange(xrange, tb, topm+height+1, 0, c.Xlabel, 0)
	TxtYRange(c.YRange, tb, leftm-2, 0, c.Ylabel, 0)

	xf := xrange.Data2Screen
	yf := c.YRange.Data2Screen
	sy0 := yf(c.YRange.Min)

	sbw := xf(2*barWidth) - xf(barWidth)

	for _, data := range c.Data {
		symbol := data.Style.Symbol
		for k, v := range data.Samples {
			x := float64(c.catIdx(k) + 1)
			y := v
			sx := xf(x - barWidth/2)
			sy := yf(y)
			sh := sy0 - sy
			tb.Block(sx, sy, sbw, sh, symbol)
		}
	}

	if kb != nil {

	}

	return tb.String()
}
