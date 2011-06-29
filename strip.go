package chart

import (
	"fmt"
	"rand"
	//	"os"
	//	"strings"
)


type StripChartData struct {
	Name string
	Data []float64
}


type StripChart struct {
	XRange        Range
	Jitter        bool
	Title, Xlabel string
	Key           Key
	Data          []StripChartData
}

func (sc *StripChart) AddData(name string, data []float64) {
	if sc.Data == nil {
		sc.Data = make([]StripChartData, 0, 1)
	}
	sc.Data = append(sc.Data, StripChartData{name, data})
	if sc.XRange.DataMin == 0 && sc.XRange.DataMax == 0 {
		sc.XRange.DataMin = data[0]
		sc.XRange.DataMax = data[0]
	}
	for _, d := range data {
		if d < sc.XRange.DataMin {
			sc.XRange.DataMin = d
		} else if d > sc.XRange.DataMax {
			sc.XRange.DataMax = d
		}
	}
	sc.XRange.Min = sc.XRange.DataMin
	sc.XRange.Max = sc.XRange.DataMax
}

func (sc *StripChart) AddIntData(name string, data []int) {
	fd := make([]float64, len(data))
	for i, d := range data {
		fd[i] = float64(d)
	}
	sc.AddData(name, fd)
}


func (sc *StripChart) PlotTxt(w, h int) string {
	if sc.Key.Pos == "" {
		sc.Key.Pos = "itr"
	}

	if h < 5 {
		h = 5
	}
	if w<10 {
		w=10
	}

	width, leftm, height, topm := w-6, 3, h-4, 1

	var kb *TextBuf
	if !sc.Key.Hide {
		maxlen := 0
		entries := make([]KeyEntry, 0)
		for s, data := range sc.Data {
			text := data.Name
			if text != "" {
				text = text[:min(len(text),w/2-7)]
				symbol := Symbol[s%len(Symbol)]
				entries = append(entries, KeyEntry{symbol,text})
				if len(text) > maxlen {
					maxlen = len(text)
				}
			}
		}
		if len(entries) > 0 {
			kh, kw := len(entries)+2, maxlen+7
			kb = NewTextBuf(kw, kh)
			if sc.Key.Border != -1 {
				kb.Rect(0, 0, kw-1, kh-1, sc.Key.Border+1, ' ')
			}
			for i, e := range entries {
				kb.Put(2, i+1, e.Symbol)
				kb.Text(5, i+1, e.Text, -1)
			}

			switch sc.Key.Pos[:2]{
			case "ol": 
				width, leftm = width-maxlen-9, leftm+kw
				sc.Key.X = 0
			case "or": 
				width = width-maxlen-9
				sc.Key.X = w - kw
			case "ot": 
				height, topm = height-kh-2, topm +kh
				sc.Key.Y = 1
			case "ob": 
				height = height-kh-2
				sc.Key.Y = topm + height + 4
			case "it":
				sc.Key.Y = topm + 1
			case "ic":
				sc.Key.Y = topm + (height-kh)/2
			case "ib":
				sc.Key.Y = topm + height - kh

			}

			switch sc.Key.Pos[:2]{
			case "ol", "or": 
				switch sc.Key.Pos[2] {
				case 't': sc.Key.Y = topm
				case 'c': sc.Key.Y = topm + (height-kh)/2
				case 'b': sc.Key.Y = topm + height - kh + 1
				}
			case "ot", "ob": 
				switch sc.Key.Pos[2] {
				case 'l': sc.Key.X = leftm
				case 'c': sc.Key.X = leftm + (width-kw)/2
				case 'r': sc.Key.X = w - kw -2
				}
			}
			if sc.Key.Pos[0] == 'i' {
				switch sc.Key.Pos[2] {
				case 'l': sc.Key.X = leftm + 1
				case 'c': sc.Key.X = leftm + (width-kw)/2
				case 'r': sc.Key.X = leftm + width - kw -1
				}

			}
		}
	}

	fmt.Printf("width=%d, height=%d, leftm=%d, topm=%d\n", width, height, leftm, topm)

	tb := NewTextBuf(w, h)


	tb.Rect(leftm, topm, width, height, 0, ' ')
	tb.Text(width/2+leftm, 0, sc.Title, 0)
	tb.Text(width/2+leftm, topm+height+2, sc.Xlabel, 0)

	var ndata = len(sc.Data)

	var ntics int
	switch {
	case width < 20:
		ntics = 2
	case width < 30:
		ntics = 3
	case width < 60:
		ntics = 4
	case width < 80:
		ntics = 5
	case width < 100:
		ntics = 7
	default:
		ntics = 10
	}
	// fmt.Printf("Requesting %d tics.\n", ntics)

	sc.XRange.Setup(ntics, width, leftm, false)

	tics := sc.XRange.Tics
	if !tics.Hide {
		for tic := tics.First; tic < tics.Last+tics.Delta/2; tic += tics.Delta {
			x := sc.XRange.Data2Screen(tic)
			lab := FmtFloat(tic)
			tb.Put(x, topm+height, '+')
			tb.Text(x,topm+height+1, lab, 0)
		}
	}

	ml := 0
	for s, data := range sc.Data {
		for _, d := range data.Data {
			x := sc.XRange.Data2Screen(d)
			y := (s+1)*height/(ndata+1) + topm
			if sc.Jitter {
				y += rand.Intn(3) - 1
			}
			tb.Put(x, y, Symbol[s%len(Symbol)])
		}
		if len(data.Name) > ml {
			ml = len(data.Name)
		}
	}

	if kb != nil {
		fmt.Printf("width=%d, height=%d, leftm=%d, topm=%d, x=%d, y=%d\n", width, height, leftm, topm, sc.Key.X, sc.Key.Y)
		tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	return tb.String()
}
