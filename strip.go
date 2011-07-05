package chart

import (
	"fmt"
	"rand"
	"math"
	//	"os"
	//	"strings"
)


type StripChartData struct {
	Name string
	Data []float64
}


type StripChart struct {
	Jitter bool
	ScatterChart
}

func (sc *StripChart) AddData(name string, data []float64) {
	n := len(sc.ScatterChart.Data) + 1
	pd := make([]Point, len(data))
	for i, d := range data {
		pd[i].X = d
		pd[i].Y = float64(n)
	}
	sc.ScatterChart.AddData(name, pd)
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
		_, _, height, topm, _, _, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)
		sc.YRange.Setup(numytics, numytics+1, height, topm, true)

		yj := math.Fabs(sc.YRange.Screen2Data(1) - sc.YRange.Screen2Data(2))
		for s, data := range sc.ScatterChart.Data {
			if data.Samples == nil {
				continue // should not happen
			}
			fmt.Printf("Set %d\n", s)
			for i, p := range data.Samples {
				r := float64(rand.Intn(3) - 1)
				fmt.Printf("r=%f, delta=%f orig=%g,%g\n", r, r*yj, p.X, p.Y)
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
