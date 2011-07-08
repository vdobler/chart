package main

import (
	"chart"
	"fmt"
	// "rand"
	// "time"
)

func main() {

	data1 := []float64{15e-7, 30e-7, 35e-7, 50e-7, 70e-7, 75e-7, 80e-7, 32e-7, 35e-7, 70e-7, 65e-7}
	data2 := []float64{10e-7, 11e-7, 12e-7, 22e-7, 25e-7, 33e-7}
	data3 := []float64{50e-7, 55e-7, 55e-7, 60e-7, 50e-7, 65e-7, 60e-7, 65e-7, 55e-7, 50e-7}

	data10 := []float64{34567, 35432, 37888, 39991, 40566, 42123, 44678}
	sc := chart.StripChart{Jitter: true}
	sc.Title = "Sample Strip Chart"
	sc.Xlabel = "x - Value"

	sc.AddData("Sample A und aaa und ccc", data1)
	sc.AddData("Sample B", data2)
	sc.AddData("Sample C", data3)

	for _, pos := range []string{"itl", "itc", "itr", "icl", "icc", "icr", "ibl", "ibc", "ibr",
		"otl", "otc", "otr", "olt", "olc", "olb", "obl", "obc", "obr", "ort", "orc", "orb"} {
		sc.Key.Pos = pos
		fmt.Printf("\nKey.Pos = %s\n", pos)
		fmt.Printf("%s\n", sc.PlotTxt(100, 30))
	}

	p := chart.ScatterChart{Title: "Sample Scatter Chart", Xlabel: "X-Value", Ylabel: "Y-Value"}
	p.AddDataPair("Sample A", data10, data1)
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	p.XRange.TicSetting.Hide, p.YRange.TicSetting.Hide = true, true
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	p.Xlabel, p.Ylabel = "", ""
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	p.XRange.TicSetting.Hide, p.YRange.TicSetting.Hide = false, false
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	pl := chart.ScatterChart{Title: "Scatter + Lines", Xlabel: "X-Value", Ylabel: "Y-Value"}
	pl.Key.Pos = "itl"
	x := []float64{-4, -3.3, -1.8, -1, 0.2, 0.8, 2, 3.1, 4, 5.3, 6, 7, 8, 9}
	y := []float64{22, 18, -3, 0, 0.5, 2, 45, 12, 16.5, 24, 30, 55, 60, 70}
	pl.AddDataPair("Measurement", x, y)
	last := len(pl.Data) - 1
	pl.Data[last].Samples[6].EX1 = 1
	pl.Data[last].Samples[6].EX2 = 1.7
	pl.Data[last].Samples[6].EY1 = 5
	pl.Data[last].Samples[6].EY2 = 12
	pl.AddFunc("Theory", func(x float64) float64 {
		if x > 5.25 && x < 5.75 {
			return 75
		}
		if x > 7.25 && x < 7.75 {
			return 500
		}
		return x * x
	})
	pl.AddLinear("Line", -4, 0, 10, 60)
	fmt.Printf("%s\n", pl.PlotTxt(100, 28))

	/*
		 steps := []int64{ 1, 5, 7, 8, 10, 30, 50, 100, 150, 300, 500, 800, 1000, 1500, 3000, 5000,8000, 10000, 15000, 20000, 30000, 50000, 70000, 100000, 200000, 400000, 800000, 1200000, 1800000, 2000000, 2200000, 2500000, 3000000, 5000000, 9000000, 2 * 9000000, 4 * 9000000 }
		 for _, step := range steps {
		 fmt.Printf("\nStep %d seconds\n", step)
		 t, v := make([]float64, 20), make([]float64, 20)
		 now := time.Seconds()
		 for i := 0; i < 20; i++ {
		 t[i] = float64(now + int64(i)*step)
		 v[i] = rand.NormFloat64() * 3
		 }
		 tl := chart.ScatterChart{Title: "Date and Time", Xlabel: "X-Value", Ylabel: "Y-Value"}
		 tl.Key.Hide = true
		 tl.XRange.Time = true
		 tl.Key.Pos = "itl"
		 tl.AddDataPair("Sample", t, v)
		 fmt.Printf("%s\n", tl.PlotTxt(100, 15))
		 }

		steps2 := []int64{10, 100, 1000, 10000, 100000, 1000000, 10000000}
		for _, step := range steps2 {
			fmt.Printf("\nStep %d seconds\n", step)
			t, v := make([]float64, 20), make([]float64, 20)
			now := time.Seconds()
			for i := 0; i < 20; i++ {
				t[i] = float64(now + int64(i)*step)
				v[i] = rand.NormFloat64() * 3
			}
			tl := chart.ScatterChart{Title: "Date and Time", Xlabel: "Numeric ", Ylabel: "Date / Time"}
			tl.Key.Hide = true
			tl.YRange.Time = true
			tl.Key.Pos = "itl"
			tl.AddDataPair("Sample", v, t)
			fmt.Printf("%s\n", tl.PlotTxt(100, 25))
		}

		hc := chart.HistChart{Title: "Histogram", Xlabel: "Value", Ylabel: "Count", ShowVal: true}
		points := make([]float64, 150)
		for i := 0; i < len(points); i++ {
			x := rand.NormFloat64()*10 + 20
			if x < 0 {
				x = 0
			} else if x > 50 {
				x = 50
			}
			points[i] = x
		}
		hc.AddData("Sample 1\nfrom today\nand yesterday ", points)
		fmt.Printf("%s\n", hc.PlotTxt(120, 20))

		points2 := make([]float64, 80)
		for i := 0; i < len(points2); i++ {
			x := rand.NormFloat64()*4 + 37
			if x < 0 {
				x = 0
			} else if x > 50 {
				x = 50
			}
			points2[i] = x
		}
		hc.AddData("Sample 2\ntomorrow", points2)
		fmt.Printf("%s\n", hc.PlotTxt(120, 20))

		hc.Stacked = true
		fmt.Printf("%s\n", hc.PlotTxt(120, 20))

		points3 := make([]float64, 60)
		for i := 0; i < len(points3); i++ {
			x := rand.NormFloat64() * 15
			if x < 0 {
				x = 0
			} else if x > 50 {
				x = 50
			}
			points3[i] = x
		}
		hc.AddData("Sample 3", points3)
		fmt.Printf("%s\n", hc.PlotTxt(120, 30))
		hc.Stacked = false
		fmt.Printf("%s\n", hc.PlotTxt(120, 30))

		hc.AddData("Sample 4\nhuhu", points3)
		hc.Stacked = true
		fmt.Printf("%s\n", hc.PlotTxt(120, 30))
		hc.Key.Cols = 2
		hc.Key.Pos = "ort"
		fmt.Printf("%s\n", hc.PlotTxt(120, 30))
		hc.Key.Cols = -3
		hc.Key.Pos = "irt"
		fmt.Printf("%s\n", hc.PlotTxt(124, 30))

		//
		// Box Charts
		//
		bc := chart.BoxChart{Title: "Box Chart", Xlabel: "Value", Ylabel: "Count"}

		for x := 10; x <= 50; x += 5 {
			p := make([]float64, 70)
			a := rand.Float64() * 10
			v := rand.Float64()*5 + 2
			for i := 0; i < len(p); i++ {
				x := rand.NormFloat64()*v + a
				p[i] = x
			}
			bc.AddSet(float64(x), p, true)
		}
		fmt.Printf("%s\n", bc.PlotTxt(100, 25))

		bc.NextDataSet("Hallo")
		for x := 12; x <= 50; x += 10 {
			p := make([]float64, 60)
			a := rand.Float64()*15 + 30
			v := rand.Float64()*5 + 2
			for i := 0; i < len(p); i++ {
				x := rand.NormFloat64()*v + a
				p[i] = x
			}
			bc.AddSet(float64(x), p, true)
		}
		fmt.Printf("%s\n", bc.PlotTxt(100, 25))
	*/

}
