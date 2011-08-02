package chart


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
type DataStyle struct {
	Symbol      int     // 0: no symbol; any codepoint: this symbol
	SymbolColor string  // 
	SymbolSize  float64 // 
	LineStyle   int     // 0: solid, 1 dashed, 2 dotted, 3 dashdotdot, 4 longdash  5 longdot
	LineColor   string  // 0: auto = same as SymbolColor
	LineWidth   int     // 0: no line
	Fill        float64 // 0: none, 1: same as line, 0.x: lighter fill
	Font        string  // "": default
	FontSize    int     // -2: tiny, -1: small, 0: normal, 1: large, 2: huge
	FontColor   string  // 
	Alpha       float64
}

const (
	NoLine = iota
	SolidLine
	DashedLine
	DottedLine
	DashDotDotLine
	LongDashLine
	LongDotLine
)

func (d *DataStyle) empty() bool {
	return d.Symbol == 0 && d.SymbolColor == "" && d.LineStyle == 0 && d.LineColor == "" && d.Fill == 0 && d.SymbolSize == 0
}


// Style is a list of suitable default styles.
var Style = []DataStyle{
	DataStyle{Symbol: 'o', SymbolColor: "#cc0000", LineStyle: 0, LineColor: "#cc0000",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
	DataStyle{Symbol: '=', SymbolColor: "#00bb00", LineStyle: 1, LineColor: "#00bb00",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
	DataStyle{Symbol: '%', SymbolColor: "#0000dd", LineStyle: 2, LineColor: "#0000dd",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
	DataStyle{Symbol: '&', SymbolColor: "#996600", LineStyle: 3, LineColor: "#996600",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
	DataStyle{Symbol: '+', SymbolColor: "#bb00bb", LineStyle: 4, LineColor: "#bb00bb",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
	DataStyle{Symbol: 'X', SymbolColor: "#00aaaa", LineStyle: 5, LineColor: "#00aaaa",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
	DataStyle{Symbol: '*', SymbolColor: "#aaaa00", LineStyle: 6, LineColor: "#aaaa00",
		Fill: 0, SymbolSize: 1, Font: "Verdana", FontSize: 0, Alpha: 0},
}


var autostylecnt int = 0

// AutoStyle produces on subsequent call new styles based on the Style list.
func AutoStyle() (style DataStyle) {
	n := len(Style)
	si := autostylecnt % n
	ci := (si + autostylecnt/n) % n
	li := (si + 2*autostylecnt/n) % n
	style.Symbol = Style[si].Symbol
	style.SymbolColor = Style[ci].SymbolColor
	style.LineColor = Style[ci].LineColor
	style.LineStyle = Style[li].LineStyle
	style.Fill = Style[autostylecnt].Fill
	style.SymbolSize = Style[autostylecnt].SymbolSize
	style.Font = Style[autostylecnt].Font
	style.FontSize = Style[autostylecnt].FontSize
	style.Alpha = Style[autostylecnt].Alpha
	autostylecnt++
	return
}
