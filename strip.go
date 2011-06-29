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
	Jitter        bool
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
	for s, data := range sc.ScatterChart.Data {
		fmt.Printf("Set: %d\n", s)
		for i, _ := range data.Data {
			fmt.Printf("%f ,%f\n", data.Data[i].X,data.Data[i].Y)
		}
	}
}


func (sc *StripChart) PlotTxt(w, h int) string {
	sc.ScatterChart.Ylabel = ""
	sc.ScatterChart.YRange.Tics.Hide = true
	sc.ScatterChart.YRange.MinMode.Fixed = true
	sc.ScatterChart.YRange.MinMode.Value = 0.5
	sc.ScatterChart.YRange.MaxMode.Fixed = true
	sc.ScatterChart.YRange.MaxMode.Value = float64(len(sc.ScatterChart.Data))+0.5

	if sc.Jitter {
		sc.LayoutTxt(w,h) // Set up ranging
		one, two := sc.YRange.Screen2Data(1), sc.YRange.Screen2Data(2)
		fmt.Printf("one: %f   two: %f\n", one, two)
		yj := math.Fabs(two-one)
		fmt.Printf("yj = %f\n", yj)
		for s, data := range sc.ScatterChart.Data {
			fmt.Printf("Set %d\n", s)
			for i, p := range data.Data {
				r := float64(rand.Intn(3)-1)
				fmt.Printf("r=%f, delta=%f orig=%g,%g\n", r, r *yj, p.X,p.Y)
				data.Data[i].Y += r * yj
			}
		}
	}
	result := sc.ScatterChart.PlotTxt(w,h)

	// Revert Jitter
	for s, data := range sc.ScatterChart.Data {
		for i, _ := range data.Data {
			data.Data[i].Y = float64(s+1)
		}
	}

	return result
}
