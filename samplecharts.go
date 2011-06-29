package main

import (
	"chart"
	"fmt"
)

func main() {
	data1 := []float64{15e-7, 30e-7, 35e-7, 50e-7, 70e-7, 75e-7, 80e-7, 32e-7, 35e-7, 70e-7, 65e-7}
	data2 := []float64{10e-7, 11e-7, 12e-7, 22e-7, 25e-7, 33e-7}
	data3 := []float64{50e-7, 55e-7, 55e-7, 60e-7, 50e-7, 65e-7}

	data10 := []float64{34567, 35432, 37888, 39991, 40566, 42123, 44678}
	sc := chart.StripChart{Title: "Sample Strip Chart", Xlabel: "X-Value", Jitter: true}

	sc.AddData("Sample A und aaa und ccc", data1)
	sc.AddData("Sample B", data2)
	sc.AddData("Sample C", data3)

	for _, pos := range []string{"itl", "itc", "itr", "icl", "icc", "icr", "ibl", "ibc", "ibr",
		"otl", "otc", "otr", "olt", "olc", "olb","obl", "obc", "obr", "ort", "orc", "orb", } {
		sc.Key.Pos = pos
		fmt.Printf("\nKey.Pos = %s\n", pos)
		fmt.Printf("%s\n", sc.PlotTxt(80, 20))
	}

	p := chart.ScatterChart{Title: "Sample Scatter Chart", Xlabel: "X-Value", Ylabel: "Y-Value"}
	p.AddDataPair("Sample A", data10, data1)
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	p.XRange.Tics.Hide, p.YRange.Tics.Hide = true, true
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	p.Xlabel, p.Ylabel = "", ""
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

	p.XRange.Tics.Hide, p.YRange.Tics.Hide = false, false
	fmt.Printf("%s\n", p.PlotTxt(80, 20))

}
