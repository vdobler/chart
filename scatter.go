package chart

import (
// "fmt"
//	"os"
//	"strings"
)


type ScatterChartData struct {
	Name    string
	Style   DataStyle
	Samples []Point
	Func    func(float64) float64
}


type ScatterChart struct {
	XRange, YRange Range
	Title          string
	Xlabel, Ylabel string
	Key            Key
	Data           []ScatterChartData
}

func (sc *ScatterChart) AddFunc(name string, f func(float64) float64) {
	s := Symbol[ len(sc.Data) % len(Symbol) ]
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, nil, f})
	sc.Key.Entries = append(sc.Key.Entries, KeyEntry{s, name})
}

func (sc *ScatterChart) AddLinear(name string, ax, ay, bx, by float64) {
	sc.AddFunc(name, func(x float64) float64 {
		return ay + (x-ax)*(by-ay)/(bx-ax)
	})
}


func (sc *ScatterChart) AddData(name string, data []Point) {
	s := Symbol[ len(sc.Data) % len(Symbol) ]
	sc.Data = append(sc.Data, ScatterChartData{name, DataStyle{}, data, nil})
	sc.Key.Entries = append(sc.Key.Entries, KeyEntry{s, name})
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
	n := min(len(x), len(y))
	data := make([]Point, n)
	for i := 0; i < n; i++ {
		data[i].X = x[i]
		data[i].Y = y[i]
	}
	sc.AddData(name, data)
}





func (sc *ScatterChart) PlotTxt(w, h int) string {
	width, leftm, height, topm, kb, numxtics, numytics := LayoutTxt(w, h, sc.Title, sc.Xlabel, sc.Ylabel, sc.XRange.TicSetting.Hide, sc.YRange.TicSetting.Hide, &sc.Key)

	sc.XRange.Setup(numxtics, numxtics+2, width, leftm, false)
	sc.YRange.Setup(numytics, numytics+1, height, topm, true)

	xlabsep, ylabsep := 1, 3
	if !sc.XRange.TicSetting.Hide {
		xlabsep++
	}
	if !sc.YRange.TicSetting.Hide {
		ylabsep += 6
	}

	tb := NewTextBuf(w, h)
	tb.Rect(leftm, topm, width, height, 0, ' ')
	if sc.Title != "" {
		tb.Text(width/2+leftm, 0, sc.Title, 0)
	}
	if sc.Xlabel != "" {
		y := topm + height + 1
		if !sc.XRange.TicSetting.Hide {
			y++
		}
		tb.Text(width/2+leftm, y, sc.Xlabel, 0)
	}
	if sc.Ylabel != "" {
		x := leftm - 3
		if !sc.YRange.TicSetting.Hide {
			x -= 6
		}
		tb.Text(x, topm+height/2, sc.Ylabel, 3)
	}

	for _, tic := range sc.XRange.Tics {
		x := sc.XRange.Data2Screen(tic.Pos)
		lx := sc.XRange.Data2Screen(tic.LabelPos)
		if sc.XRange.Time {
			tb.Put(x, topm+height, '|')
			tb.Put(x, topm+height+1, '|')
			if tic.Align == -1 {
				tb.Text(lx+1, topm+height+1, tic.Label, -1)
			} else {
				tb.Text(lx, topm+height+1, tic.Label, 0)
			}
		} else {
			tb.Put(x, topm+height, '+')
			tb.Text(lx, topm+height+1, tic.Label, 0)
		}
		if sc.XRange.Time {
			tb.Text(leftm, topm+height+2, sc.XRange.TMin.Format("2006-01-02 15:04:05"), -1)
			tb.Text(leftm+width+1, topm+height+2, sc.XRange.TMax.Format("2006-01-02 15:04:05"), 1)
		}
	}

	for _, tic := range sc.YRange.Tics {
		y := sc.YRange.Data2Screen(tic.Pos)
		ly := sc.YRange.Data2Screen(tic.LabelPos)
		tb.Put(leftm, y, '+')
		tb.Text(leftm-1, ly, tic.Label, 1)
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
			for sx := leftm; sx < leftm+width; sx++ {
				x := sc.XRange.Screen2Data(sx)
				y := data.Func(x)
				sy := sc.YRange.Data2Screen(y)
				if y >= sc.YRange.Min && y <= sc.YRange.Max {
					tb.Put(sx, sy, symbol)
				}
				d := abs(lastsy - sy)
				// fmt.Printf("Point (%.2f, %.2f) : sx=%d, sy=%d\n", x, y, sx, sy)
				if sx > leftm && d > 2 {
					// Oversampling
					f := 1
					if sy < lastsy {
						f = -1
					}
					osx := lastsx
					// fmt.Printf("Oversampling: d=%d, f=%d, from %d to %d\n", d, f, lastsy+f, sy-f)
					var done bool
					for osy := clip(lastsy+f, 0, h); osy != clip(sy-f, 0, h); osy += f {
						// fmt.Printf("  osx=%d, osy=%d\n", osx, osy)
						if sc.YRange.Screen2Data(osy) >= sc.YRange.Min && sc.YRange.Screen2Data(osy) <= sc.YRange.Max {
							tb.Put(osx, osy, symbol)
						}
						if !done && osy > (sy+lastsy)/2 {
							osx++
							done = true
						}
					}
				}

				lastsx, lastsy = sx, sy
			}
		}
	}

	if kb != nil {
		tb.Paste(sc.Key.X, sc.Key.Y, kb)
	}

	return tb.String()
}
