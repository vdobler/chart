package chart

import (
	"fmt"
	"math"
	//	"os"
	//	"strings"
)

type CatValue struct {
	Cat string
	Val float64
}

type CategoryChartData struct {
	Name    string
	Style   DataStyle
	Samples []CatValue
}


type PieChart struct {
	Title   string
	Key     Key
	ShowVal bool // Display values 
	InnerR  float64
	Data    []CategoryChartData
}

func (c *PieChart) AddData(name string, data []CatValue) {
	s := Symbol[len(c.Data)%len(Symbol)]
	c.Data = append(c.Data, CategoryChartData{name, DataStyle{Symbol: s}, data})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{s, name})
}

func (c *PieChart) AddDataPair(name string, cat []string, val []float64) {
	n := min(len(cat), len(val))
	data := make([]CatValue, n)
	for i := 0; i < n; i++ {
		data[i].Cat, data[i].Val = cat[i], val[i]
	}
	c.AddData(name, data)
}


func (c *PieChart) PlotTxt(w, h int) string {
	tb := NewTextBuf(w, h)

	// TODO(vodo): handel w<h case
	left, top := 2, 1
	radiusy := (h - top*2) / 2
	if c.Title != "" {
		radiusy -= 1
		top += 1
	}
	radiusx := 2 * radiusy

	x0, y0 := left+radiusx, h/2

	fmt.Printf("w,h = %d,%d;   left,top=%d,%d;    ry,rx=%d,%d;   x0,y0=%d,%d\n",
		w, h, left, top, radiusy, radiusx, x0, y0)
	if c.Title != "" {
		tb.Text(left+radiusx, 0, c.Title, 0)
	}

	for i, data := range c.Data {
		var sum float64
		for _, d := range data.Samples {
			sum += d.Val
		}
		fmt.Printf("Total of set %d: %f\n", i, sum)

		var phi float64 = -math.Pi
		for j, d := range data.Samples {
			symbol := Symbol[j%len(Symbol)]
			alpha := 2 * math.Pi * d.Val / sum
			for r := c.InnerR * float64(radiusy); r <= float64(radiusy)+0.1; r += 0.2 {
				for w := phi; w < phi+alpha-0.01; w += 0.02 {
					x, y := int(0.5+2*r*math.Cos(w)), int(0.5+r*math.Sin(w))
					tb.Put(x+x0, y+y0, symbol)
				}
			}
			phi += alpha
		}
	}

	// if kb != nil {
	//	tb.Paste(c.Key.X, c.Key.Y, kb)
	//}

	return tb.String()
}
