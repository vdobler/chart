package chart

import (
	"fmt"
	"math"
	//	"os"
	//	"strings"
)


type HistChartData struct {
	Name      string
	DataStyle DataStyle
	Samples   []float64
}


type HistChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Horizontal     bool // Display is horizontal bars
	Stacked        bool // Display different data sets ontop of each other
	ShowVal        bool
	Data           []HistChartData
	FirstBin       float64 // center of the first (lowest bin)
	BinWidth       float64
	TBinWidth      TimeDelta // for time XRange
}

func (c *HistChart) AddData(name string, data []float64) {
	s := Symbol[len(c.Data)%len(Symbol)]
	c.Data = append(c.Data, HistChartData{name, DataStyle{}, data})
	c.Key.Entries = append(c.Key.Entries, KeyEntry{s, name})

	if len(c.Data) == 1 {  // first data set 
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
}

func (c *HistChart) AddDataInt(name string, data []int) {
	fdata := make([]float64, len(data))
	for i, d:=range data {
		fdata[i] = float64(d)
	}
	c.AddData(name, fdata)
}

func (c *HistChart) AddDataGeneric(name string, data []Value) {
	fdata := make([]float64, len(data))
	for i, d:=range data {
		fdata[i] = d.XVal()
	}
	c.AddData(name, fdata)
}

func (hc *HistChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, hc.Title, hc.Xlabel, hc.Ylabel, hc.XRange.TicSetting.Hide, hc.YRange.TicSetting.Hide, &hc.Key)

	// Outside bound ranges for histograms are nicer
	leftm, width = leftm+2, width-2
	topm, height = topm, height-1

	hc.XRange.Setup(numxtics, numxtics+1, width, leftm, false)
	hc.BinWidth = hc.XRange.TicSetting.Delta
	binCnt := int((hc.XRange.Max-hc.XRange.Min)/hc.BinWidth + 0.5)
	hc.FirstBin = hc.XRange.Min + hc.BinWidth/2

	counts := make([][]int, len(hc.Data))
	hc.YRange.DataMin = 0
	max := 0
	for i, data := range hc.Data {
		count := make([]int, binCnt)
		for _, x := range data.Samples {
			bin := int((x - hc.XRange.Min) / hc.BinWidth)
			count[bin] = count[bin] + 1
			if count[bin] > max {
				max = count[bin]
			}
		}
		counts[i] = count
		// fmt.Printf("Count: %v\n", count)
	}
	if hc.Stacked { // recalculate max
		max = 0
		for bin := 0; bin < binCnt; bin++ {
			sum := 0
			for i := range counts {
				sum += counts[i][bin]
			}
			// fmt.Printf("sum of bin %d = %d\n", bin, sum)
			if sum > max {
				max = sum
			}
		}
	}
	hc.YRange.DataMax = float64(max)
	hc.YRange.Setup(numytics, numytics+2, height, topm, true)

	tb := NewTextBuf(w, h)

	if hc.Title != "" {
		tb.Text(width/2+leftm, 0, hc.Title, 0)
	}

	TxtXRange(hc.XRange, tb, topm+height+1, 0, hc.Xlabel, 0)
	TxtYRange(hc.YRange, tb, leftm-2, 0, hc.Ylabel, 0)

	xf := hc.XRange.Data2Screen
	yf := hc.YRange.Data2Screen

	numSets := len(hc.Data)
	for i, tic := range hc.XRange.Tics {
		xs := xf(tic.Pos)
		// tb.Put(xs, topm+height+1, '+')
		// tb.Text(lx, topm+height+2, tic.Label, 0)

		if i == 0 {
			continue
		}

		last := hc.XRange.Tics[i-1]
		lasts := xf(last.Pos)

		var blockW int
		if hc.Stacked {
			blockW = xs - lasts - 1
		} else {
			blockW = int(float64(xs-lasts-numSets) / float64(numSets))
		}
		// fmt.Printf("blockW= %d\n", blockW)

		center := (tic.Pos + last.Pos) / 2
		bin := int((center - hc.XRange.Min) / hc.BinWidth)
		xs = lasts
		lastCnt := 0
		y0 := yf(0)

		minCnt := int(math.Fabs(hc.YRange.Screen2Data(0)-hc.YRange.Screen2Data(1)) / 2)

		for d, _ := range hc.Data {
			cnt := counts[d][bin]
			y := yf(float64(lastCnt + cnt))
			if cnt > minCnt {
				fill := Symbol[d%len(Symbol)]

				tb.Block(xs+1, y, blockW, y0-y, fill)

				if hc.ShowVal {
					lab := fmt.Sprintf("%d", cnt)
					if blockW-len(lab) >= 4 {
						lab = " " + lab + " "
					}
					xlab := xs + blockW/2 + 1 // hc.XRange.Data2Screen(center)
					if blockW%2 == 1 {
						xlab++
					}
					ylab := y - 1
					if numSets > 1 {
						ylab = yf(float64(lastCnt) + float64(cnt)/2)
					}
					tb.Text(xlab, ylab, lab, 0)
					// fmt.Printf("Set %d: %s at %d\n", d, lab, ylab)
				}
			}
			if !hc.Stacked {
				xs += blockW + 1
			} else {
				lastCnt += cnt
				y0 = y
			}
		}
	}

	if kb != nil {
		tb.Paste(hc.Key.X, hc.Key.Y, kb)
	}

	return tb.String()
}
