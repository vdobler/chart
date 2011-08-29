package chart

import (
	"fmt"
	"math"
	//	"os"
	// "strings"
)


// PieChart represents pie and ring charts.
type PieChart struct {
	Title string              // The title
	Key   Key                 // The Key/Legend
	Inner float64             // relative radius of inner white are (set to 0.7 to produce ring chart)
	Data  []CategoryChartData // The data

	FmtVal func(value, sume float64) string // Produce labels
}


// AbsoluteValue will format value (ignoring sum).  It is a convenience function which
// can be assigned to PieChart.FmtVal.
func AbsoluteValue(value, sum float64) (s string) {
	fv := math.Fabs(value)
	switch {
	case fv < 0.01:
		s = fmt.Sprintf(" %g ", value)
	case fv < 0.1:
		s = fmt.Sprintf(" %.2f ", value)
	case fv < 1:
		s = fmt.Sprintf(" %.1f ", value)
	case fv < 100000:
		s = fmt.Sprintf(" %.0f ", value)
	default:
		s = fmt.Sprintf(" %g ", value)
	}
	return
}

// PercentValue formats value as percentage of sum.   It is a convenience function which
// can be assigned to PieChart.FmtVal.
func PercentValue(value, sum float64) (s string) {
	value *= 100 / sum
	s = AbsoluteValue(value, sum) + "% "
	return
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


var PieChartShrinkage = 0.66 // Scaling factor of radius of next data set.
var PieChartHighlight = 0.15 // How much are flaged segments offset. 

// Plot outputs the scatter chart sc to g.
func (c *PieChart) Plot(g Graphics) {
	layout := layout(g, c.Title, "", "", true, true, &c.Key)

	width, height := layout.Width, layout.Height
	topm, leftm := layout.Top, layout.Left
	width += 0

	r := imin(height, width) / 2
	x0, y0 := leftm+r, topm+r

	// Make sure pie fits into plotting area
	rshift := int(float64(r) * PieChartHighlight)
	if rshift < 6 {
		rshift = 6
	}
	for _, d := range c.Data[0].Samples {
		if d.Flag {
			// debug.Printf("Reduced %d by %d", r, rshift)
			r -= rshift / 3
			break
		}
	}

	g.Begin()

	if c.Title != "" {
		g.Title(c.Title)
	}

	for _, data := range c.Data {

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
			shift := 0

			var t string
			if c.FmtVal != nil {
				t = c.FmtVal(d.Val, sum)
			}
			if d.Flag {
				shift = rshift
			}

			wedges[j] = Wedgeinfo{Phi: phi, Psi: phi + alpha, Text: t, Tp: "c",
				Style: style, Font: Font{}, Shift: shift}

			phi += alpha
		}
		g.Rings(wedges, x0, y0, r, ri)

		r = int(float64(r) * PieChartShrinkage)
	}

	if !c.Key.Hide {
		g.Key(layout.KeyX, layout.KeyY, c.Key)
	}

	g.End()
}
