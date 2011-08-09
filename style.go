package chart

import (
	"math"
	"fmt"
)

// Symbol is a list of different symbols. 
var Symbol = []int{'o', // empty circle
	'=', // empty square
	'%', // empty triangle up
	'&', // empty diamond
	'+', // plus
	'X', // cross
	'*', // star
	'0', // bulls eys
	'@', // filled circle
	'#', // filled square
	'A', // filled tringale up
	'Z', // filled diamond
	'.', // tiny dot
}

// SymbolIndex returns the index of the symbol s in Symbol or -1 if not found.
func SymbolIndex(s int) (idx int) {
	for idx = 0; idx < len(Symbol); idx++ {
		if Symbol[idx] == s {
			return idx
		}
	}
	return -1
}

// NextSymbol returns the next symbol of s: Either in the global list Symbol 
// or (if not found there) the next character.
func NextSymbol(s int) int {
	if idx := SymbolIndex(s); idx != -1 {
		return Symbol[(idx+1)%len(Symbol)]
	}
	return s + 1
}

var CharacterWidth = map[int]float32{'a': 16.8, 'b': 17.0, 'c': 15.2, 'd': 16.8, 'e': 16.8, 'f': 8.5, 'g': 17.0,
	'h': 16.8, 'i': 5.9, 'j': 5.9, 'k': 16.8, 'l': 6.9, 'm': 25.5, 'n': 16.8, 'o': 16.8, 'p': 17.0, 'q': 17.0,
	'r': 10.2, 's': 15.2, 't': 8.4, 'u': 16.8, 'v': 15.4, 'w': 22.2, 'x': 15.2, 'y': 15.2, 'z': 15.2,
	'A': 20.2, 'B': 20.2, 'C': 22.2, 'D': 22.2, 'E': 20.2, 'F': 18.6, 'G': 23.5, 'H': 22.0, 'I': 8.2, 'J': 15.2,
	'K': 20.2, 'L': 16.8, 'M': 25.5, 'N': 22.0, 'O': 23.5, 'P': 20.2, 'Q': 23.5, 'R': 21.1, 'S': 20.2, 'T': 18.5,
	'U': 22.0, 'V': 20.2, 'W': 29.0, 'X': 20.2, 'Y': 20.2, 'Z': 18.8, ' ': 8.5,
	'1': 16.8, '2': 16.8, '3': 16.8, '4': 16.8, '5': 16.8, '6': 16.8, '7': 16.8, '8': 16.8, '9': 16.8, '0': 16.8,
	'.': 8.2, ',': 8.2, ':': 8.2, ';': 8.2, '+': 17.9, '"': 11.0, '*': 11.8, '%': 27.0, '&': 20.2, '/': 8.4,
	'(': 10.2, ')': 10.2, '=': 18.0, '?': 16.8, '!': 8.5, '[': 8.2, ']': 8.2, '{': 10.2, '}': 10.2, '$': 16.8,
	'<': 18.0, '>': 18.0, '§': 16.8, '°': 12.2, '^': 14.2, '~': 18.0,
}
var averageCharacterWidth float32

func init() {
	n := 0
	for _, w := range CharacterWidth {
		averageCharacterWidth += w
		n++
	}
	averageCharacterWidth /= float32(n)
	averageCharacterWidth = 15
}

var Palette = map[string]string{"title": "#aa9933", "label": "#000000", "axis": "#000000",
	"ticlabel": "#000000", "grid": "#c0c0c0", "keyborder": "#000000", "errorbar": "*0.3",
}

// DataStyle contains all information about all graphic elements in a chart.
// TODOs:
//  - remove Font..., not part of DataStyle and relevant only for Text
//  - keep Symbol as "show this symbol in strip/scatter/box plots"
//  - add new Char as "char/symbol to use as text replacement for color"
//    that would be for "lines without symbols", hist, bar, cbar, pie
//    box
//
// differentiate between drawing data/plot-style in scatter (points, lines, linespoints)
// and style (color, symbol, width, filling). disalow e.g. in "datastyle lines"
// linwidth of 0.
//
type DataStyle struct {
	Symbol      int     // 0: no symbol; any codepoint: this symbol
	SymbolColor string  // 
	SymbolSize  float64 // 
	LineStyle   int     // 0: solid, 1 dashed, 2 dotted, 3 dashdotdot, 4 longdash  5 longdot
	LineColor   string  // 0: auto = same as SymbolColor
	LineWidth   int     // 0: no line
	FillColor   string  // "": no fill
	Alpha       float64 // 
}

// PlotStyle describes how data and functions are drawn in scatter plots.
// Can be used to describe how a key entry is drawn
type PlotStyle int

const (
	PlotStylePoints      = 1
	PlotStyleLines       = 2
	PlotStyleLinesPoints = 3
	PlotStyleBox         = 4
)

func (ps PlotStyle) undefined() bool {
	return int(ps) < 1 || int(ps) > 3
}


const (
	SolidLine = iota
	DashedLine
	DottedLine
	DashDotDotLine
	LongDashLine
	LongDotLine
)

type Font struct {
	Name  string // "": default
	Size  int    // -2: tiny, -1: small, 0: normal, 1: large, 2: huge
	Color string // "": default, other: use this
}

func (d *DataStyle) empty() bool {
	return d.Symbol == 0 && d.SymbolColor == "" && d.LineStyle == 0 && d.LineColor == "" && d.FillColor == "" && d.SymbolSize == 0
}


// Standard colors used by AutoStyle
var StandardColors = []string{"#cc0000", "#00bb00", "#0000dd", "#996600", "#bb00bb", "#00aaaa",
	"#aaaa00"}
// Standard line styles used by AutoStyle (fill=false)
var StandardLineStyles = []int{SolidLine, DashedLine, DottedLine, LongDashLine, LongDotLine}
// Standard symbols used by AutoStyle
var StandardSymbols = []int{'o', '=', '%', '&', '+', 'X', '*', '@', '#', 'A', 'Z'}
// How much brighter/darker filled elements become.
var StandardFillFactor = 0.5


// AutoStyle produces a styles based on StandardColors, StandardLineStyles, and StandardSymbols.
// Call with fill = true for charts with filled elements (hist, bar, cbar, pie).
func AutoStyle(i int, fill bool) (style DataStyle) {
	nc, nl, ns := len(StandardColors), len(StandardLineStyles), len(StandardSymbols)

	si := i % ns
	ci := i % nc
	li := i % nl

	style.Symbol = StandardSymbols[si]
	style.SymbolColor = StandardColors[ci]
	style.LineColor = StandardColors[ci]
	style.SymbolSize = 1
	style.Alpha = 0

	if fill {
		style.LineStyle = SolidLine
		style.LineWidth = 2
		if i < nc {
			style.FillColor = lighter(style.LineColor, StandardFillFactor)
		} else if i <= 2*nc {
			style.FillColor = darker(style.LineColor, StandardFillFactor)
		} else {
			style.FillColor = style.LineColor
		}
	} else {
		style.LineStyle = StandardLineStyles[li]
		style.LineWidth = 1
	}
	return
}


var DefaultStyle = map[string]DataStyle{"axis": DataStyle{LineColor: "#000000", LineWidth: 2, LineStyle: SolidLine},
	"maxis": DataStyle{LineColor: "#000000", LineWidth: 2, LineStyle: SolidLine}, // mirrored axis
	"tic":   DataStyle{LineColor: "#000000", LineWidth: 1, LineStyle: SolidLine},
	"mtic":  DataStyle{LineColor: "#000000", LineWidth: 1, LineStyle: SolidLine},
	"zero":  DataStyle{LineColor: "#404040", LineWidth: 1, LineStyle: SolidLine},
	"grid":  DataStyle{LineColor: "#808080", LineWidth: 1, LineStyle: SolidLine},
	"key":   DataStyle{LineColor: "#202020", LineWidth: 1, LineStyle: SolidLine, FillColor: "#f0f0f0", Alpha: 0.2},
}

var DefaultFont = map[string]Font{"title": Font{Size: +1}, "label": Font{}, "key": Font{Size: -1},
	"tic": Font{}, "rangelimit": Font{},
}

func hsv2rgb(h, s, v int) (r, g, b int) {
	H := int(math.Floor(float64(h) / 60))
	S, V := float64(s)/100, float64(v)/100
	f := float64(h)/60 - float64(H)
	p := V * (1 - S)
	q := V * (1 - S*f)
	t := V * (1 - S*(1-f))

	switch H {
	case 0, 6:
		r, g, b = int(255*V), int(255*t), int(255*p)
	case 1:
		r, g, b = int(255*q), int(255*V), int(255*p)
	case 2:
		r, g, b = int(255*p), int(255*V), int(255*t)
	case 3:
		r, g, b = int(255*p), int(255*q), int(255*V)
	case 4:
		r, g, b = int(255*t), int(255*p), int(255*V)
	case 5:
		r, g, b = int(255*V), int(255*p), int(255*q)
	default:
		panic(fmt.Sprintf("Ooops: Strange H value %d in hsv2rgb(%d,%d,%d).", H, h, s, v))
	}

	return
}

func f3max(a, b, c float64) float64 {
	switch true {
	case a > b && a >= c:
		return a
	case b > c && b >= a:
		return b
	case c > a && c >= b:
		return c
	}
	return a
}

func f3min(a, b, c float64) float64 {
	switch true {
	case a < b && a <= c:
		return a
	case b < c && b <= a:
		return b
	case c < a && c <= b:
		return c
	}
	return a
}

func rgb2hsv(r, g, b int) (h, s, v int) {
	R, G, B := float64(r)/255, float64(g)/255, float64(b)/255

	if R == G && G == B {
		h, s = 0, 0
		v = int(r * 255)
	} else {
		max, min := f3max(R, G, B), f3min(R, G, B)
		if max == R {
			h = int(60 * (G - B) / (max - min))
		} else if max == G {
			h = int(60 * (2 + (B-R)/(max-min)))
		} else {
			h = int(60 * (4 + (R-G)/(max-min)))
		}
		if max == 0 {
			s = 0
		} else {
			s = int(100 * (max - min) / max)
		}
		v = int(100 * max)
	}
	if h < 0 {
		h += 360
	}
	return
}

func color2rgb(color string) (r, g, b int) {
	if color[0] == '#' {
		color = color[1:]
	}
	n, err := fmt.Sscanf(color, "%2x%2x%2x", &r, &g, &b)
	if n != 3 || err != nil {
		r, g, b = 127, 127, 127
	}
	// fmt.Printf("%s  -->  %d %d %d\n", color,r,g,b)
	return
}


func lighter(color string, f float64) string {
	r, g, b := color2rgb(color)
	h, s, v := rgb2hsv(r, g, b)
	f = 1 - f
	s = int(float64(s) * f)
	v += int((100 - float64(v)) * f)
	if v > 100 {
		v = 100
	}
	r, g, b = hsv2rgb(h, s, v)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func darker(color string, f float64) string {
	r, g, b := color2rgb(color)
	h, s, v := rgb2hsv(r, g, b)
	f = 1 - f
	v = int(float64(v) * f)
	s += int((100 - float64(s)) * f)
	if s > 100 {
		s = 100
	}
	r, g, b = hsv2rgb(h, s, v)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
