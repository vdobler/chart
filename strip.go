package chart

import (
	// "fmt"
	"rand"
	"math"
	//	"os"
	//	"strings"
)


type StripChart struct {
	Jitter bool
	ScatterChart
}

func (sc *StripChart) AddData(name string, data []float64, style DataStyle) {
	n := len(sc.ScatterChart.Data) + 1
	pd := make([]EPoint, len(data))
	nan := math.NaN()
	for i, d := range data {
		pd[i].X = d
		pd[i].Y = float64(n)
		pd[i].DeltaX, pd[i].DeltaY = nan, nan
	}
	if style.empty() {
		style = AutoStyle()
	}
	style.LineStyle = 0
	sc.ScatterChart.AddData(name, pd, style)
}

func (sc *StripChart) AddDataGeneric(name string, data []Value) {
	n := len(sc.ScatterChart.Data) + 1
	pd := make([]EPoint, len(data))
	nan := math.NaN()
	for i, d := range data {
		pd[i].X = d.XVal()
		pd[i].Y = float64(n)
		pd[i].DeltaX, pd[i].DeltaY = nan, nan
	}
	sc.ScatterChart.AddData(name, pd, DataStyle{})
}


func (sc *StripChart) PlotTxt(w, h int) string {
	sc.ScatterChart.Ylabel = ""
	sc.ScatterChart.YRange.TicSetting.Hide = true
	sc.ScatterChart.YRange.MinMode.Fixed = true
	sc.ScatterChart.YRange.MinMode.Value = 0.5
	sc.ScatterChart.YRange.MaxMode.Fixed = true
	sc.ScatterChart.YRange.MaxMode.Value = float64(len(sc.ScatterChart.Data)) + 0.5

	if sc.Jitter {
		// Set up ranging
		_, _, height, topm, _, _, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key, 1, 1)
		sc.YRange.Setup(numytics, numytics+1, height, topm, true)

		yj := math.Fabs(sc.YRange.Screen2Data(1) - sc.YRange.Screen2Data(2))
		for _, data := range sc.ScatterChart.Data {
			if data.Samples == nil {
				continue // should not happen
			}
			for i := range data.Samples {
				r := float64(rand.Intn(3) - 1)
				data.Samples[i].Y += r * yj
			}
		}
	}
	result := sc.ScatterChart.PlotTxt(w, h)

	// Revert Jitter
	for s, data := range sc.ScatterChart.Data {
		if data.Samples == nil {
			continue // should not happen
		}
		for i, _ := range data.Samples {
			data.Samples[i].Y = float64(s + 1)
		}
	}

	return result
}


func (sc *StripChart) Plot(g Graphics) {
	sc.ScatterChart.Ylabel = ""
	sc.ScatterChart.YRange.TicSetting.Hide = true
	sc.ScatterChart.YRange.TicSetting.Delta = 1
	sc.ScatterChart.YRange.MinMode.Fixed = true
	sc.ScatterChart.YRange.MinMode.Value = 0.5
	sc.ScatterChart.YRange.MaxMode.Fixed = true
	sc.ScatterChart.YRange.MaxMode.Value = float64(len(sc.ScatterChart.Data)) + 0.5
	
	if sc.Jitter {
		// Set up ranging
		layout := Layout(g, sc.Title, sc.XRange.Label, sc.YRange.Label,
			sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)

		_, height := layout.Width, layout.Height
		topm, _ := layout.Top, layout.Left
		_, numytics := layout.NumXtics, layout.NumYtics
		
		sc.YRange.Setup(numytics, numytics+1, height, topm, true)
		
		// amplitude of jitter: not too smal to be visible and useful, not to
		// big to be ugly or even overlapp other
		
		null := sc.YRange.Screen2Data(0)
		absmin := 1.4 * math.Fabs(sc.YRange.Screen2Data(1) - null ) // would be one pixel
		tenpc := math.Fabs(sc.YRange.Screen2Data(height) - null)/10 // 10 percent of graph area
		smplcnt := len(sc.ScatterChart.Data) + 1 //  as samples are borders
		noverlp := math.Fabs(sc.YRange.Screen2Data(height/smplcnt) - null) // do not overlapp other sample
			
		yj := noverlp
		if tenpc < yj {
			yj = tenpc
		}
		if yj < absmin { 
			yj = absmin 
		}
		
		// yjs := sc.YRange.Data2Screen(yj) - sc.YRange.Data2Screen(0)
		// fmt.Printf("yj = %.2f : in screen = %d\n", yj, yjs)
		for _, data := range sc.ScatterChart.Data {
			if data.Samples == nil {
				continue // should not happen
			}
			for i := range data.Samples {
				shift := yj * rand.NormFloat64() * yj
				data.Samples[i].Y += shift
			}
		}
	}
	sc.ScatterChart.Plot(g)

	if sc.Jitter {
		// Revert Jitter
		for s, data := range sc.ScatterChart.Data {
			if data.Samples == nil {
				continue // should not happen
			}
			for i, _ := range data.Samples {
				data.Samples[i].Y = float64(s + 1)
			}
		}
	}
}
