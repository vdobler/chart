package chart

import (
	// "fmt"
	//	"os"
	//	"strings"
)


type ScatterChartData struct {
	Name string
	Style DataStyle
	Samples []Point
	Func func(float64)float64
}


type ScatterChart struct {
	XRange, YRange   Range
	Title  string
	Xlabel, Ylabel string
	Key           Key
	Data          []ScatterChartData
}

func (sc *ScatterChart) AddFunc(name string, f func(float64)float64) {
	if sc.Data == nil {
		sc.Data = make([]ScatterChartData, 0, 1)
	}
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, nil, f})
}

func (sc *ScatterChart) AddLinear(name string, ax, ay, bx, by float64) {
	sc.AddFunc(name, func(x float64)float64 {
		return ay + (x-ax)*(by-ay)/(bx-ax)
	})
}


func (sc *ScatterChart) AddData(name string, data []Point) {
	if sc.Data == nil {
		sc.Data = make([]ScatterChartData, 0, 1)
	}
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, data, nil})
	if sc.XRange.DataMin == 0 && sc.XRange.DataMax == 0 && sc.YRange.DataMin == 0 && sc.YRange.DataMax == 0 {
		sc.XRange.DataMin = data[0].X
		sc.XRange.DataMax = data[0].X
		sc.YRange.DataMin = data[0].Y
		sc.YRange.DataMax = data[0].Y
	}
	for _, d := range data {
		if d.X < sc.XRange.DataMin {
			sc.XRange.DataMin = d.X
		} else if d.X > sc.XRange.DataMax {
			sc.XRange.DataMax = d.X
		}
		if d.Y < sc.YRange.DataMin {
			sc.YRange.DataMin = d.Y
		} else if d.Y > sc.YRange.DataMax {
			sc.YRange.DataMax = d.Y
		}
	}
	sc.XRange.Min = sc.XRange.DataMin
	sc.XRange.Max = sc.XRange.DataMax
	sc.YRange.Min = sc.YRange.DataMin
	sc.YRange.Max = sc.YRange.DataMax
	// fmt.Printf("New Limits: x %f %f; y %f %f\n", sc.XRange.DataMin, sc.XRange.DataMax, sc.YRange.DataMin, sc.YRange.DataMax) 
}

func (sc *ScatterChart) AddDataPair(name string, x, y []float64) {
	n := min(len(x),len(y))
	data := make([]Point, n)
	for i:=0; i<n; i++ {
		data[i].X = x[i]
		data[i].Y = y[i]
	}
	sc.AddData(name, data)
}


func (sc *ScatterChart) LayoutTxt(w, h int) (width, leftm, height, topm int, kb *TextBuf) {
	if sc.Key.Pos == "" {
		sc.Key.Pos = "itr"
	}

	if h < 5 {
		h = 5
	}
	if w<10 {
		w=10
	}

	width, leftm, height, topm = w-4, 2, h-1, 0
	xlabsep, ylabsep := 1, 3
	if sc.Title != "" {
		topm++
		height--
	}
	if sc.Xlabel != "" {
		height--
	}
	if !sc.XRange.Tics.Hide {
		height--
		xlabsep++
	}
	if sc.Ylabel != "" {
		leftm += 2
		width -= 2
	}
	if !sc.YRange.Tics.Hide {
		leftm += 6
		width -= 6
		ylabsep += 6
	}


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

	// fmt.Printf("width=%d, height=%d, leftm=%d, topm=%d\n", width, height, leftm, topm)

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
	// fmt.Printf("Requesting %d,%d tics.\n", ntics,height/3)

	sc.XRange.Setup(ntics, width, leftm, false)
	sc.YRange.Setup(height/3, height, topm, true)
	
	return
} 


func (sc *ScatterChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb := sc.LayoutTxt(w,h)

	xlabsep, ylabsep := 1, 3
	if !sc.XRange.Tics.Hide {
		xlabsep++
	}
	if !sc.YRange.Tics.Hide {
		ylabsep += 6
	}


	tb := NewTextBuf(w, h)
	tb.Rect(leftm, topm, width, height, 0, ' ')
	if sc.Title != "" {
		tb.Text(width/2+leftm, 0, sc.Title, 0)
	}
	if sc.Xlabel != "" {
		y :=  topm + height + 1
		if !sc.XRange.Tics.Hide { y++ }
		tb.Text(width/2+leftm, y, sc.Xlabel, 0)
	}
	if sc.Ylabel != "" {
		x := leftm - 3
		if !sc.YRange.Tics.Hide { x -= 6 }
		tb.Text(x, topm+height/2, sc.Ylabel, 3)
	}

	tics := sc.XRange.Tics
	if !tics.Hide {
		for tic := tics.First; tic < tics.Last+tics.Delta/2; tic += tics.Delta {
			x := sc.XRange.Data2Screen(tic)
			lab := FmtFloat(tic)
			tb.Put(x, topm+height, '+')
			tb.Text(x,topm+height+1, lab, 0)
		}
	}

	tics = sc.YRange.Tics
	if !tics.Hide {
		for tic := tics.First; tic < tics.Last+tics.Delta/2; tic += tics.Delta {
			y := sc.YRange.Data2Screen(tic)
			lab := FmtFloat(tic)
			tb.Put(leftm, y, '+')
			tb.Text(leftm-1 ,y , lab, 1)
		}
	}

	// Plot Data
	for s, data := range sc.Data {
		if data.Samples != nil {
			// Samples
			for _, d := range data.Samples {
				x := sc.XRange.Data2Screen(d.X)
				y := sc.YRange.Data2Screen(d.Y)
				tb.Put(x, y, Symbol[s%len(Symbol)])
			}
		} else if data.Func != nil {
			// Functions
			var lastsy, lastsx int
			symbol := Symbol[s%len(Symbol)]
			for sx:=leftm; sx<leftm+width; sx++ {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				sy := sc.YRange.Data2Screen(y)
				if y>=sc.YRange.Min && y<=sc.YRange.Max {
					tb.Put(sx, sy, symbol)
				}
				d := abs(lastsy-sy)
				// fmt.Printf("Point (%.2f, %.2f) : sx=%d, sy=%d\n", x, y, sx, sy)
				if sx > leftm && d>2 {
					// Oversampling
					f := 1
					if sy < lastsy { f = -1 }
					osx := lastsx
					// fmt.Printf("Oversampling: d=%d, f=%d, from %d to %d\n", d, f, lastsy+f, sy-f)
					var done bool
					for osy:=clip(lastsy+f,0,h); osy!=clip(sy-f,0,h); osy+=f {
						// fmt.Printf("  osx=%d, osy=%d\n", osx, osy)
						if sc.YRange.Screen2Data(osy) >= sc.YRange.Min && sc.YRange.Screen2Data(osy)<=sc.YRange.Max {
							tb.Put(osx, osy, symbol)
						}
						if !done && osy > (sy+lastsy)/2 {
							osx++
							done = true 
						}
					}
				}
				
				lastsx, lastsy  = sx, sy
			}
		}
	}

	if kb != nil {
		//fmt.Printf("width=%d, height=%d, leftm=%d, topm=%d, x=%d, y=%d\n", width, 
		//	height, leftm, topm, sc.Key.X, sc.Key.Y)
		tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	return tb.String()
}
