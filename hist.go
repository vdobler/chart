package chart

import (
	"fmt"
//	"os"
//	"strings"
)


type HistChartData struct {
	Name string
	DataStyle DataStyle
	Samples []float64
}


type HistChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Horizontal     bool  // Display is horizontal bars
	Stacked        bool  // Display different data sets ontop of each other
	Data           []HistChartData
	FirstBin       float64  // center of the first (lowest bin)
	BinWidth       float64
	TBinWidth      TimeDelta // for time XRange
}

func (c *HistChart) AddData(name string, data []float64) {
	if c.Data == nil {
		c.Data = make([]HistChartData, 0, 1)
	}
	c.Data = append(c.Data, HistChartData{name, DataStyle{}, data})
	if c.XRange.DataMin == 0 && c.XRange.DataMax == 0  {
		c.XRange.DataMin = data[0]
		c.XRange.DataMax = data[0]
	}
	for _, d := range data {
		if d < c.XRange.DataMin {
			c.XRange.DataMin = d
		} else if d > c.XRange.DataMax {
			c.XRange.DataMax = d
		}
	}
	c.XRange.Min = c.XRange.DataMin
	c.XRange.Max = c.XRange.DataMax
	// fmt.Printf("New Limits: x %f %f; y %f %f\n", 
	// 	c.XRange.DataMin, c.XRange.DataMax, c.YRange.DataMin, c.YRange.DataMax) 
}


func (hc *HistChart) PlotTxt(w, h int) string {
	width, height, leftm, topm := w - 10, h - 10, 5, 5
	ntics := 10
	hc.XRange.Setup(ntics, ntics+5, width, leftm, false)
	hc.BinWidth = hc.XRange.TicSetting.Delta
	binCnt := int((hc.XRange.Max - hc.XRange.Min) / hc.BinWidth  + 0.5)
	hc.FirstBin = hc.XRange.Min + hc.BinWidth/2

	count := make([]int, binCnt)
	hc.YRange.DataMin = 0
	max := 0
	for _, x := range hc.Data[0].Samples {
		bin := int((x - hc.XRange.Min)/hc.BinWidth)
		count[bin] = count[bin] + 1
		if count[bin] > max {
			max = count[bin]
		}
	}
	hc.YRange.DataMax = float64(max)
	hc.YRange.Setup(height/5, height/5+3, height, topm, true)

	tb := NewTextBuf(w, h)
	tb.Rect(leftm, topm, width, height, 0, ' ')
	if hc.Title != "" {
		tb.Text(width/2+leftm, 0, hc.Title, 0)
	}

	for i, tic := range hc.XRange.Tics {
		x := hc.XRange.Data2Screen(tic.Pos)
		lx := hc.XRange.Data2Screen(tic.LabelPos)
		tb.Put(x, topm+height, '+')
		tb.Text(lx, topm+height+1, tic.Label, 0)

		y0 := hc.YRange.Data2Screen(0)
		if i > 0 {
			last := hc.XRange.Tics[i-1]
			center := (tic.Pos + last.Pos)/2
			bin := int((center - hc.XRange.Min) / hc.BinWidth)
			cnt := count[bin]
			y := hc.YRange.Data2Screen(float64(cnt))
			x0 := hc.XRange.Data2Screen(last.Pos)
			// fmt.Printf("Bin x=%d, y=%d, w=%d, h=%d\n",x0, y, x-x0, y-y0)
			tb.Rect(x0, y, x-x0, y0-y, 0, '#')
			lab := fmt.Sprintf("%d", cnt)
			xlab := hc.XRange.Data2Screen(center)
			y--
			tb.Text(xlab, y, lab, 0 )
			
		}
	}

	for _, tic := range hc.YRange.Tics {
		y := hc.YRange.Data2Screen(tic.Pos)
		ly := hc.YRange.Data2Screen(tic.LabelPos)
		tb.Put(leftm, y, '+')
		tb.Text(leftm-1, ly, tic.Label, 1)
	}


	return tb.String()
}
