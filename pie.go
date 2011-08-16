package chart

import (
	"fmt"
	"math"
	//	"os"
	// "strings"
)


// PieChart represents pie and ring charts.
type PieChart struct {
	Title   string
	Key     Key
	ShowVal int     // Display values. 0: don't show, 1: relative in %, 2: absolute 
	Inner   float64 // relative radius of inner white are (set to 0.7 to produce ring chart)
	Data    []CategoryChartData
}

type CategoryChartData struct {
	Name    string
	Style   []Style
	Samples []CatValue
}


func (c *PieChart) AddData(name string, data []CatValue, style []Style) {
	if len(style) < len(data) {
		ns := make([]Style, len(data))
		copy(style, ns)
		for i := len(style); i < len(data); i++ {
			ns[i] = AutoStyle(i-len(style), true)
		}
		style = ns
	}
	c.Data = append(c.Data, CategoryChartData{name, style, data})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{PlotStyle: -1, Text: name})
	for s, cv := range data {
		c.Key.Entries = append(c.Key.Entries, KeyEntry{PlotStyle: PlotStyleBox, Style: style[s], Text: cv.Cat})
	}
}

func (c *PieChart) AddDataPair(name string, cat []string, val []float64) {
	n := imin(len(cat), len(val))
	data := make([]CatValue, n)
	for i := 0; i < n; i++ {
		data[i].Cat, data[i].Val = cat[i], val[i]
	}
	c.AddData(name, data, nil)
}


func (c *PieChart) formatVal(v, sum float64) (s string) {
	if c.ShowVal == 1 {
		v *= 100 / sum // percentage
	}
	switch {
	case v < 0.1:
		s = fmt.Sprintf(" %.2f ", v)
	case v < 1:
		s = fmt.Sprintf(" %.1f ", v)
	default:
		s = fmt.Sprintf(" %.0f ", v)
	}
	if c.ShowVal == 1 {
		s += "% "
	}
	return
}

var PieChartTextAscpect float64 = 1.9 // how much wider is the x-radius
var PieChartShrinkage = 0.65          // Scaling factor of radius of next data set.
var PieChartBorder = 0.05             // Fraction of white border outside next data sets.


// Plot outputs the scatter chart sc to g.
func (c *PieChart) Plot(g Graphics) {
	layout := Layout(g, c.Title, "", "", true, true, &c.Key)

	width, height := layout.Width, layout.Height
	topm, leftm := layout.Top, layout.Left
	width += 0

	r := height / 2
	r0 := r
	x0, y0 := leftm+r, topm+r

	g.Begin()

	if c.Title != "" {
		g.Title(c.Title)
	}

	for i, data := range c.Data {
		// _ := c.Key.Entries[keidx].Text // data set name
		// style = DataStyle{LineColor: "#404040", LineWidth: 3, LineStyle: SolidLine}

		var sum float64
		for _, d := range data.Samples {
			sum += d.Val
		}

		wedges := make([]Wedgeinfo, len(data.Samples))
		var ri int = 0
		if c.Inner > 0 {
			ri = int(float64(r) * c.Inner)
		}

		var phi float64 = -math.Pi
		for j, d := range data.Samples {
			style := data.Style[j]
			alpha := 2 * math.Pi * d.Val / sum

			var t string
			if c.ShowVal > 0 {
				t = c.formatVal(d.Val, sum)
			}
			wedges[j] = Wedgeinfo{Ro: r, Ri: ri, Phi: phi, Psi: phi + alpha, Text: t, Tp: "c", Style: style, Font: Font{}}

			phi += alpha
		}
		g.Rings(wedges, x0, y0, r0)

		r = int(float64(r) * PieChartShrinkage)
		if i < len(c.Data)-1 {
			ra := int(float64(r) * (1 + PieChartBorder))
			if ra == r {
				ra++
			}
			st := Style{LineWidth: 0, FillColor: "#ffffff", Symbol: ' '}
			g.Wedge(x0, y0, ra, 0, 15, st)
		}
	}

	if !c.Key.Hide {
		g.Key(layout.KeyX, layout.KeyY, c.Key)
	}

	g.End()
}
