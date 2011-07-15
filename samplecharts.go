package main

import (
	"chart"
	"fmt"
	"os"
	// "rand"
	// "time"
)

func main() {

	data1 := []float64{15e-7, 30e-7, 35e-7, 50e-7, 70e-7, 75e-7, 80e-7, 32e-7, 35e-7, 70e-7, 65e-7}
	data10 := []float64{34567, 35432, 37888, 39991, 40566, 42123, 44678}
	/*
		data2 := []float64{10e-7, 11e-7, 12e-7, 22e-7, 25e-7, 33e-7}
		data3 := []float64{50e-7, 55e-7, 55e-7, 60e-7, 50e-7, 65e-7, 60e-7, 65e-7, 55e-7, 50e-7}

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
	*/

	p := chart.ScatterChart{Title: "Sample Scatter Chart", Xlabel: "X-Value", Ylabel: "Y-Value"}
	p.AddDataPair("Sample A", data10, data1, chart.DataStyle{})
	fmt.Printf("%s\n", p.PlotTxt(100, 25))

	/*
		p.XRange.TicSetting.Hide, p.YRange.TicSetting.Hide = true, true
		fmt.Printf("%s\n", p.PlotTxt(100, 25))

		p.Xlabel, p.Ylabel = "", ""
		fmt.Printf("%s\n", p.PlotTxt(100, 25))

		p.XRange.TicSetting.Hide, p.YRange.TicSetting.Hide = false, false
		fmt.Printf("%s\n", p.PlotTxt(100, 25))
	*/

	pl := chart.ScatterChart{Title: "Scatter + Lines", Xlabel: "X-Value", Ylabel: "Y-Value"}
	pl.Key.Pos = "itl"
	// pl.XRange.TicSetting.Delta = 5
	pl.XRange.TicSetting.Grid = 1
	x := []float64{-4, -3.3, -1.8, -1, 0.2, 0.8, 2, 3.1, 4, 5.3, 6, 7, 8, 9}
	y := []float64{22, 18, -3, 0, 0.5, 2, 45, 12, 16.5, 24, 30, 55, 60, 70}
	pl.AddDataPair("Measurement", x, y, chart.AutoStyle())
	last := len(pl.Data) - 1
	pl.Data[last].Samples[6].DeltaX = 2.5
	pl.Data[last].Samples[6].OffX = 0.5
	pl.Data[last].Samples[6].DeltaY = 16
	pl.Data[last].Samples[6].OffY = 2
	pl.AddData("Volker", []chart.EPoint{chart.EPoint{-4, 40, 0, 0, 0, 0}, chart.EPoint{-3, 45, 0, 0, 0, 0},
		chart.EPoint{-2, 35, 0, 0, 0, 0}},
		chart.AutoStyle())
	pl.AddFunc("Theory", func(x float64) float64 {
		if x > 5.25 && x < 5.75 {
			return 75
		}
		if x > 7.25 && x < 7.75 {
			return 500
		}
		return x * x
	},chart.DataStyle{})
	fmt.Printf("%s\n", pl.PlotTxt(100, 28))

	sf, _ := os.Create("scatter.svg")
	pl.PlotSvg(600, 400, sf)
	sf.Close()

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

	// Bar chart
	barc := chart.BarChart{Title: "My first Bar Chart"}
	barc.XRange.ShowZero = true
	barc.AddDataPair("Amount",
		[]float64{-10, 10, 20, 30, 35, 40, 50},
		[]float64{90, 120, 180, 205, 230, 150, 190})
	fmt.Printf("%s\n", barc.PlotTxt(100, 25))
	barc.AddDataPair("Test",
		[]float64{-5, 15, 25, 35, 45, 55},
		[]float64{110, 80, 95, 80, 120, 140})
	fmt.Printf("%s\n", barc.PlotTxt(100, 25))
	barc.SameBarWidth = true
	fmt.Printf("%s\n", barc.PlotTxt(100, 25))

	// Pie Chart
	piec := chart.PieChart{Title: "Some Pies"}
	piec.AddDataPair("Europe", []string{"D", "AT", "CH", "F", "E", "I"}, []float64{10, 20, 30, 35, 15, 25})
	piec.Inner = 0.5
	piec.ShowVal = 1
	fmt.Printf("%s\n", piec.PlotTxt(80, 30))
	piec.AddDataPair("America", []string{"North", "Middel", "South"}, []float64{20, 10, 15})
	piec.Key.Cols = 2
	fmt.Printf("%s\n", piec.PlotTxt(80, 30))

	// Categorized Bar Chart
	cbarc := chart.CategoryBarChart{Title: "Income", Categories: []string{"none", "low", "average", "high"}}
	cbarc.AddData("Europe", map[string]float64{"none": 10, "low": 15, "average": 25, "high": 20})
	fmt.Printf("%s\n", cbarc.PlotTxt(100, 20))
	cbarc.AddData("Asia", map[string]float64{"none": 15, "low": 30, "average": 10, "high": 20})
	fmt.Printf("%s\n", cbarc.PlotTxt(100, 20))
	cbarc.Stacked = true
	fmt.Printf("%s\n", cbarc.PlotTxt(100, 20))

	cbarc = chart.CategoryBarChart{Title: "Income", Categories: []string{"none", "low", "average", "high"}}
	cbarc.YRange.ShowZero = true
	cbarc.AddData("Europe", map[string]float64{"none": 10, "low": 15, "average": 25, "high": 20})
	cbarc.AddData("Asia", map[string]float64{"none": 15, "low": 30, "average": 10, "high": -20})
	fmt.Printf("%s\n", cbarc.PlotTxt(100, 25))

	// Log-X axis
	lc := chart.ScatterChart{Title: "Log/Lin", Xlabel: "X-Value", Ylabel: "Y-Value"}
	lc.Key.Hide = true
	lc.XRange.Log, lc.YRange.Log = true, true
	lx := []float64{4e-2, 3e-1, 2e0, 1e1, 8e1, 7e2, 5e3}
	ly := []float64{10, 30, 90, 270, 3 * 270, 9 * 270, 27 * 270}
	lc.AddDataPair("Measurement", lx, ly, chart.DataStyle{Symbol: 'Z', SymbolColor: "#9966ff", Size: 1.5})
	fmt.Printf("%s\n", lc.PlotTxt(100, 25))

	svgf, _ := os.Create("first.svg")
	lc.PlotSvg(600, 400, svgf)
	svgf.Close()
}
