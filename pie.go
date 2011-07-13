package chart

import (
	"fmt"
	"math"
	//	"os"
	"strings"
)


type CategoryChartData struct {
	Name    string
	Style   DataStyle
	Samples []CatValue
}


type PieChart struct {
	Title   string
	Key     Key
	ShowVal int     // Display values. 0: don't show, 1: relative in %, 2: absolute 
	Inner   float64 // relative radius of inner white are (set to 0.7 to produce ring chart)
	Data    []CategoryChartData
}

func (c *PieChart) AddData(name string, data []CatValue) {
	c.Data = append(c.Data, CategoryChartData{name, DataStyle{}, data})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{Symbol: -1, Text: name})
	for s, cv := range data {
		symbol := Symbol[s%len(Symbol)]
		c.Key.Entries = append(c.Key.Entries, KeyEntry{Symbol: symbol, Text: cv.Cat})
	}
}

func (c *PieChart) AddDataPair(name string, cat []string, val []float64) {
	n := min(len(cat), len(val))
	data := make([]CatValue, n)
	for i := 0; i < n; i++ {
		data[i].Cat, data[i].Val = cat[i], val[i]
	}
	c.AddData(name, data)
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
var PieChartLabelPos = 0.75           // relativ to outer radius
var PieChartShrinkage = 0.65          // fractrion of next data set

func (c *PieChart) PlotTxt(w, h int) string {
	tb := NewTextBuf(w, h)

	// TODO(vodo): handel w<h case
	left, top := 2, 1
	radiusy := float64(h-top*2) / 2
	if c.Title != "" {
		radiusy -= 1
		top += 1
	}
	radiusx := PieChartTextAscpect * radiusy
	x0, y0 := left+int(radiusx+0.5), h/2

	dalpha := 1 / (1.5 * radiusx)

	fmt.Printf("w,h = %d,%d;   left,top=%d,%d;    ry,rx=%d,%d;   x0,y0=%d,%d\n",
		w, h, left, top, radiusy, radiusx, x0, y0)
	if c.Title != "" {
		tb.Text(left+int(radiusx), 0, c.Title, 0)
	}

	keidx := 0 // key-entry-index
	for i, data := range c.Data {
		// _ := c.Key.Entries[keidx].Text // data set name
		keidx++

		var sum float64
		for _, d := range data.Samples {
			sum += d.Val
		}

		var phi float64 = -math.Pi
		for _, d := range data.Samples {
			symbol := c.Key.Entries[keidx].Symbol
			keidx++
			alpha := 2 * math.Pi * d.Val / sum
			for r := c.Inner * radiusy; r <= radiusy+0.1; r += 0.2 {
				for w := phi + dalpha; w < phi+alpha-dalpha; w += dalpha / 5 {
					x, y := int(0.5+PieChartTextAscpect*r*math.Cos(w)), int(0.5+r*math.Sin(w))
					tb.Put(x+x0, y+y0, symbol)
				}
			}
			if i > 0 { // clear a border
				r := radiusy
				for w := float64(0); w <= 2*math.Pi; w += dalpha / 5 {
					x, y := int(0.5+PieChartTextAscpect*r*math.Cos(w)), int(0.5+r*math.Sin(w))
					tb.Put(x+x0, y+y0, ' ')
				}
			}
			if c.ShowVal != 0 {
				t := c.formatVal(d.Val, sum)
				ry, rx := PieChartLabelPos*radiusy, PieChartLabelPos*radiusx
				w := phi + alpha/2
				x, y := int(0.5+rx*math.Cos(w)), int(0.5+ry*math.Sin(w))
				tb.Text(x+x0, y+y0, t, 0)
				if radiusy > 9 {
					tb.Text(x+x0, y+y0-1, strings.Repeat(" ", len(t)), 0)
					tb.Text(x+x0, y+y0+1, strings.Repeat(" ", len(t)), 0)
				}
			}
			phi += alpha
		}
		radiusx, radiusy = radiusx*PieChartShrinkage, radiusy*PieChartShrinkage // next data set is smaler
	}

	// TODO(vodo) honour key placement
	kb := c.Key.LayoutKeyTxt()
	if kb != nil {
		tb.Paste(w-kb.W-1, 2, kb)
	}

	return tb.String()
}
