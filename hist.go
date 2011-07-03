package chart

import (
// "fmt"
//	"os"
//	"strings"
)


type HistChartData struct {
	ChartData
	Bins []Bin
}


type HistChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Data           []HistChartData
}


func (hc *HistChart) PlotTxt(w, h int) string {
	return "Oooops"
}
