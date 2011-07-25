package main

import (
	"chart"
	"fmt"
	"os"
	"math"
	"github.com/ajstarks/svgo"
	"rand"
	// "time"
)

var (
	data1 = []float64{15e-7, 30e-7, 35e-7, 50e-7, 70e-7, 75e-7, 80e-7, 32e-7, 35e-7, 70e-7, 65e-7}
	data10 = []float64{34567, 35432, 37888, 39991, 40566, 42123, 44678}

	data2 = []float64{10e-7, 11e-7, 12e-7, 22e-7, 25e-7, 33e-7}
	data3 = []float64{50e-7, 55e-7, 55e-7, 60e-7, 50e-7, 65e-7, 60e-7, 65e-7, 55e-7, 50e-7}
)



//
// Some sample strip charts
//
func stripChart() {
	file, _ := os.Create("xstrip1.svg")
	thesvg := svg.New(file)
	thesvg.Start(800, 600)
	thesvg.Title("Srip Chart")
	thesvg.Rect(0, 0, 800, 600, "fill: #ffffff")
	svggraphics := chart.NewSvgGraphics(thesvg, 400, 300, "Arial", 12)
	txtgraphics := chart.NewTextGraphics(80,25)

	c := chart.StripChart{}

	c.AddData("Sample A", data1, chart.DataStyle{})
	c.AddData("Sample B", data2, chart.DataStyle{})
	c.AddData("Sample C", data3, chart.DataStyle{})

	c.Title = "Sample Strip Chart (no Jitter)"
	c.XRange.Label = "X - Axis"
	c.Key.Pos = "icr"
	c.Plot(svggraphics)
	c.Plot(txtgraphics)
	fmt.Printf("%s\n", txtgraphics.String())

	thesvg.Gtransform("translate(400 0)")
	c.Jitter = true
	c.Title = "Sample Strip Chart (with Jitter)"
	c.Plot(svggraphics)
	c.Plot(txtgraphics)
	fmt.Printf("%s\n", txtgraphics.String())
	thesvg.Gend()

	thesvg.Gtransform("translate(0 300)")
	c.Key.Hide = true
	c.Plot(svggraphics)
	c.Plot(txtgraphics)
	fmt.Printf("%s\n", txtgraphics.String())
	thesvg.Gend()

	thesvg.Gtransform("translate(400 300)")
	c.Jitter = false
	c.Title = "Sample Strip Chart (no Jitter)"
	c.Plot(svggraphics)
	c.Plot(txtgraphics)
	fmt.Printf("%s\n", txtgraphics.String())
	thesvg.Gend()

	thesvg.End()
	file.Close()
}


//
// All different key styles
// 
func keyStyles() {
	file, _ := os.Create("xkey1.svg")
	thesvg := svg.New(file)
	w, h := 400, 300
	nw, nh := 6,6
	thesvg.Start(nw*w, nh*h)
	thesvg.Title("Key Placements")
	thesvg.Rect(0, 0, nw*w, nh*h, "fill: #ffffff")

	svggraphics := chart.NewSvgGraphics(thesvg, w, h, "Arial", 10)
	p := chart.ScatterChart{Title: "Key Placement"}
	p.XRange.TicSetting.Mirror, p.YRange.TicSetting.Mirror = 1, 1
	p.XRange.MinMode.Fixed, p.XRange.MaxMode.Fixed = true, true
	p.XRange.MinMode.Value, p.XRange.MaxMode.Value = -5, 5
	p.XRange.Min, p.XRange.Max = -5, 5
	p.XRange.TicSetting.Delta = 2

	p.YRange.MinMode.Fixed, p.YRange.MaxMode.Fixed = true, true
	p.YRange.MinMode.Value, p.YRange.MaxMode.Value = -5, 5
	p.YRange.Min, p.YRange.Max = -5, 5
	p.YRange.TicSetting.Delta = 3

	p.AddFunc("Sin", func(x float64) float64 { return math.Sin(x) }, chart.DataStyle{LineColor: "#a00000", LineWidth: 1, LineStyle: 1})
	p.AddFunc("Cos", func(x float64) float64 { return math.Cos(x) }, chart.DataStyle{LineColor: "#00a000", LineWidth: 1, LineStyle: 1})
	p.AddFunc("Tan", func(x float64) float64 { return math.Tan(x) }, chart.DataStyle{LineColor: "#0000a0", LineWidth: 1, LineStyle: 1})

	x, y := 0, 0
	for _, pos := range []string{"itl", "itc", "itr", "icl", "icc", "icr", "ibl", "ibc", "ibr",
		"otl", "otc", "otr", "olt", "olc", "olb", "obl", "obc", "obr", "ort", "orc", "orb"} {
		p.Key.Pos = pos
		p.Title = "Key Placement: " + pos
		thesvg.Gtransform(fmt.Sprintf("translate(%d %d)", x, y))
		p.Plot(svggraphics)
		thesvg.Gend()
		
		x += w
		if x+w > nw*w {
			x, y = 0, y + h
		}
	 }

	p.Key.Pos = "itl"
	p.AddFunc("Log", func(x float64) float64 { return math.Log(x) }, chart.DataStyle{LineColor: "#ff6060", LineWidth: 1, LineStyle: 1})
	p.AddFunc("Exp", func(x float64) float64 { return math.Exp(x) }, chart.DataStyle{LineColor: "#60ff60", LineWidth: 1, LineStyle: 1})
	p.AddFunc("Atan", func(x float64) float64 { return math.Atan(x) }, chart.DataStyle{LineColor: "#6060ff", LineWidth: 1, LineStyle: 1})
	p.AddFunc("Y1", func(x float64) float64 { return math.Y1(x) }, chart.DataStyle{LineColor: "#d0d000", LineWidth: 1, LineStyle: 1})

	for _, cols := range []int{ -4,-3,-2,-1,0,1,2,3,4} {
		p.Key.Cols = cols
		p.Title = fmt.Sprintf("Key Cols: %d", cols)
		thesvg.Gtransform(fmt.Sprintf("translate(%d %d)", x, y))
		p.Plot(svggraphics)
		thesvg.Gend()
		
		x += w
		if x+w > nw*w {
			x, y = 0, y + h
		}
	}

	thesvg.End()
	file.Close()
}	


//
// Scatter plots with different tic/grid settings
//
func scatterTics() {
	file, _ := os.Create("xscatter1.svg")
	thesvg := svg.New(file)
	thesvg.Start(800, 600)
	thesvg.Title("Srip Chart")
	thesvg.Rect(0, 0, 800, 600, "fill: #ffffff")
	svggraphics := chart.NewSvgGraphics(thesvg, 400, 300, "Arial", 12)

	p := chart.ScatterChart{Title: "Sample Scatter Chart"}
	p.AddDataPair("Sample A", data10, data1, chart.DataStyle{})
	p.XRange.TicSetting.Delta = 5000
	p.XRange.Label = "X - Value"
	p.YRange.Label = "Y - Value"

	p.Plot(svggraphics)

	thesvg.Gtransform("translate(400 0)")
	p.XRange.TicSetting.Hide, p.YRange.TicSetting.Hide = true, true
	p.Plot(svggraphics)
	thesvg.Gend()

	thesvg.Gtransform("translate(0 300)")
	p.YRange.TicSetting.Hide = false
	p.XRange.TicSetting.Grid, p.YRange.TicSetting.Grid = 1, 1
	p.Plot(svggraphics)
	thesvg.Gend()

	thesvg.Gtransform("translate(400 300)")
	p.XRange.TicSetting.Hide, p.YRange.TicSetting.Hide = false, false
	p.XRange.TicSetting.Mirror, p.YRange.TicSetting.Mirror = 1, 2
	p.Plot(svggraphics)
	thesvg.Gend()

	thesvg.End()
	file.Close()
}


//
// Full fletched scatter plots
//
func fancyScatter() {
	pl := chart.ScatterChart{Title: "Scatter + Lines", Xlabel: "X-Value", Ylabel: "Y-Value"}
	pl.Key.Pos = "itl"
	// pl.XRange.TicSetting.Delta = 5
	pl.XRange.TicSetting.Grid = 1
	x := []float64{-4, -3.3, -1.8, -1, 0.2, 0.8, 1.8, 3.1, 4, 5.3, 6, 7, 8, 9}
	y := []float64{22, 18, -3, 0, 0.5, 2, 45, 12, 16.5, 24, 30, 55, 60, 70}
	pl.AddDataPair("Mmnt", x, y, chart.DataStyle{Symbol: '#', SymbolColor: "#0000ff", LineStyle: 0})
	last := len(pl.Data) - 1
	pl.Data[last].Samples[6].DeltaX = 2.5
	pl.Data[last].Samples[6].OffX = 0.5
	pl.Data[last].Samples[6].DeltaY = 16
	pl.Data[last].Samples[6].OffY = 2
	pl.AddData("abcde", []chart.EPoint{chart.EPoint{-4, 40, 0, 0, 0, 0}, chart.EPoint{-3, 45, 0, 0, 0, 0},
		chart.EPoint{-2, 35, 0, 0, 0, 0}},
		chart.DataStyle{Symbol: '0', SymbolColor: "#ff00ff", LineStyle: 1, LineWidth: 1})
	pl.AddFunc("wxyz", func(x float64) float64 {
		if x > 5.25 && x < 5.75 {
			return 75
		}
		if x > 7.25 && x < 7.75 {
			return 500
		}
		return x * x
	},chart.DataStyle{Symbol: 0, LineWidth: 2, LineColor: "#a00000", LineStyle: 1})
	pl.AddFunc("30", func(x float64) float64 { return 30 },
		chart.DataStyle{Symbol: 0, LineWidth: 1, LineColor: "#00a000", LineStyle: 1})
	pl.AddFunc("7777", func(x float64) float64 { return 7 },
		chart.DataStyle{Symbol: 0, LineWidth: 1, LineColor: "#0000a0", LineStyle: 1})


	pl.XRange.ShowZero = true
	pl.XRange.TicSetting.Mirror = 1
	pl.YRange.TicSetting.Mirror = 1
	pl.XRange.TicSetting.Grid = 1
	pl.XRange.Label = "X-Range"
	pl.YRange.Label = "Y-Range"
	pl.Key.Cols = 2
	pl.Key.Pos = "orb"

	s2f, _ := os.Create("xscatter2.svg")
	mysvg := svg.New(s2f)
	mysvg.Start(1000, 600)
	mysvg.Title("My Plot")
	mysvg.Rect(0, 0, 1000, 600, "fill: #ffffff")
	svggraphics := chart.NewSvgGraphics(mysvg, 1000, 600, "Arial", 18)
	txtgraphics := chart.NewTextGraphics(100,30)
	pl.Plot(svggraphics)
	pl.Plot(txtgraphics)
	fmt.Printf("%s\n", txtgraphics.String())
	mysvg.End()
	s2f.Close()
}



//
// Box Charts
//
func boxChart() {
	file, _ := os.Create("xbox1.svg")
	thesvg := svg.New(file)
	thesvg.Start(800, 600)
	thesvg.Title("Srip Chart")
	thesvg.Rect(0, 0, 800, 600, "fill: #ffffff")
	svggraphics := chart.NewSvgGraphics(thesvg, 400, 300, "Arial", 12)

	p := chart.BoxChart{Title: "Box Chart"}
	p.XRange.Label, p.YRange.Label = "Value", "Count"

	for x := 10; x <= 50; x += 5 {
		points := make([]float64, 70)
		a := rand.Float64() * 10
		v := rand.Float64()*5 + 2
		for i := 0; i < len(points); i++ {
			x := rand.NormFloat64()*v + a
			points[i] = x
		}
		p.AddSet(float64(x), points, true)
	}

	p.NextDataSet("Hallo", chart.DataStyle{LineColor: "#00c000", LineWidth: 1, LineStyle: chart.SolidLine})
	for x := 12; x <= 50; x += 10 {
		points := make([]float64, 60)
		a := rand.Float64()*15 + 30
		v := rand.Float64()*5 + 2
		for i := 0; i < len(points); i++ {
			x := rand.NormFloat64()*v + a
			points[i] = x
		}
		p.AddSet(float64(x), points, true)
	}
	
	p.Plot(svggraphics)

	thesvg.End()
	file.Close()
}

func main() {
	stripChart()

	scatterTics()

	keyStyles()

	boxChart()
	
	fancyScatter()

	goto ende


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
	lc.AddDataPair("Measurement", lx, ly, chart.DataStyle{Symbol: 'Z', SymbolColor: "#9966ff", SymbolSize: 1.5})
	fmt.Printf("%s\n", lc.PlotTxt(100, 25))

ende:

	if true {
		s2f, _ := os.Create("text.svg")
		mysvg := svg.New(s2f)
		mysvg.Start(1600, 800)
		mysvg.Title("My Plot")
		mysvg.Rect(0, 0, 2000, 800, "fill: #ffffff")
		sgr := chart.NewSvgGraphics(mysvg, 2000, 800, "Arial", 18)
		sgr.Begin()

		texts := []string{"ill", "WWW", "Some normal text.", "Illi, is. illigalli: ill!", "OO WORKSHOOPS OMWWW BMWWMB"}
		fonts := []string{"Arial", "Helvetica", "Times", "Courier" /* "Calibri", "Palatino" */ }
		sizes := []int{8, 12, 16, 20}
		style := chart.DataStyle{FontColor: "#000000", Alpha: 0, LineColor: "#ff0000", LineWidth: 2}
		ls := chart.DataStyle{FontColor: "#000a0", Font: "Arial", FontSize: 12, Alpha: 0, LineColor: "#000000", LineWidth: 1}

		x, y := 20, 40
		for _, t := range texts {
			for _, f := range fonts {
				for _, s := range sizes {
					style.Font, style.FontSize = f, s
					tvl := sgr.TextLen(t, style)
					sgr.Text(x+tvl/2, y-2, t, "cc", 0, style)
					sgr.Line(x, y, x+tvl, y, style)
					sgr.Text(x+tvl+10, y-2, f, "cl", 0, ls)
					y += 30
					if y > 760 {
						y = 40
						x += 300
					}
				}
			}
		}

		sgr.End()
		mysvg.End()
		s2f.Close()

	}

}
