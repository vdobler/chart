package chart

import (
// "fmt"
//	"os"
//	"strings"
)


type BoxChartData struct {
	Name    string
	Style   DataStyle
	Samples []Box
}


type BoxChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Data           []BoxChartData
}


// Will add to last dataset one new box calculated from data.
// If outlier than outliers (1.5*IQR from 25/75 percentil) are
// drawn, else the wiskers extend from min to max.
func (c *BoxChart) AddSet(x float64, data []float64, outlier bool) {
	min, lq, med, avg, uq, max := SixvalFloat64(data, 25)
	b := Box{X: x, Avg: avg, Med: med, Q1: lq, Q3: uq, Low: min, High: max}

	if len(c.Data) == 0 {
		c.Data = make([]BoxChartData, 1)
		c.Data[0] = BoxChartData{Name: "", Style: DataStyle{}}
	}

	if len(c.Data) == 1 && len(c.Data[0].Samples) == 0 {
		c.XRange.DataMin, c.XRange.DataMax = x, x
		c.YRange.DataMin, c.YRange.DataMax = min, max
	} else {
		if x < c.XRange.DataMin {
			c.XRange.DataMin = x
		} else if x > c.XRange.DataMax {
			c.XRange.DataMax = x
		}
		if min < c.YRange.DataMin {
			c.YRange.DataMin = min
		}
		if max > c.YRange.DataMax {
			c.YRange.DataMax = max
		}
	}

	if outlier {
		outliers := make([]float64, 0)
		iqr := uq - lq
		min, max = max, min
		for _, d := range data {
			if d > uq+1.5*iqr || d < lq-1.5*iqr {
				outliers = append(outliers, d)
			}
			if d > max && d <= uq+1.5*iqr {
				max = d
			}
			if d < min && d >= lq-1.5*iqr {
				min = d
			}
		}
		b.Low, b.High, b.Outliers = min, max, outliers
	}
	j := len(c.Data) - 1
	c.Data[j].Samples = append(c.Data[j].Samples, b)
}


func (c *BoxChart) NextDataSet(name string) {
	s := Symbol[len(c.Data)%len(Symbol)]
	c.Data = append(c.Data, BoxChartData{name, DataStyle{}, nil})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{s, name})
}


func (c *BoxChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, c.Title, c.Xlabel, c.Ylabel, c.XRange.TicSetting.Hide, c.YRange.TicSetting.Hide, &c.Key)

	c.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	c.YRange.Setup(numytics, numytics+1, height, topm, true)

	xlabsep, ylabsep := 1, 3
	if !c.XRange.TicSetting.Hide {
		xlabsep++
	}
	if !c.YRange.TicSetting.Hide {
		ylabsep += 6
	}

	tb := NewTextBuf(w, h)
	// tb.Rect(leftm, topm, width, height, 0, ' ')
	if c.Title != "" {
		tb.Text(width/2+leftm, 0, c.Title, 0)
	}

	TxtXRange(c.XRange, tb, topm+height, topm, c.Xlabel, 2)
	TxtYRange(c.YRange, tb, leftm, leftm+width, c.Ylabel, 2)

	// Plot Data
	for s, data := range c.Data {
		// Samples
		hbw := 2 // Half Box Width
		nums := len(data.Samples)
		mhw := width / (2*nums - 1)
		if mhw > 7 {
			hbw = 3
		} else if mhw < 5 {
			hbw = 1
		}

		symbol := Symbol[s%len(Symbol)]

		for _, d := range data.Samples {
			x := c.XRange.Data2Screen(d.X)
			low, high := c.YRange.Data2Screen(d.Low), c.YRange.Data2Screen(d.High)
			q1, q3 := c.YRange.Data2Screen(d.Q1), c.YRange.Data2Screen(d.Q3)
			avg, med := c.YRange.Data2Screen(d.Avg), c.YRange.Data2Screen(d.Med)

			tb.Rect(x-hbw, q1, 2*hbw, q3-q1, 0, ' ')
			tb.Put(x-hbw, med, '+')
			for i := 0; i < hbw; i++ {
				tb.Put(x-i, med, '-')
				tb.Put(x+i, med, '-')
			}
			tb.Put(x+hbw, med, '+')

			tb.Put(x, avg, symbol)
			for y := high; y < q3; y++ {
				tb.Put(x, y, '|')
			}
			for y := q1 + 1; y <= low; y++ {
				tb.Put(x, y, '|')
			}

			for _, ol := range d.Outliers {
				y := c.YRange.Data2Screen(ol)
				tb.Put(x, y, 'o')
			}

			// tb.Put(x, y, Symbol[s%len(Symbol)])
		}
	}

	if kb != nil {
		tb.Paste(c.Key.X, c.Key.Y, kb)
	}

	return tb.String()
}
