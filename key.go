package chart

import (
	"strings"
	// "fmt"
)

// Key encapsulates settings for keys/legends in a chart.
//
// Key placement os governed by Pos which may take the following values:
//          otl  otc  otr      
//         +-------------+ 
//     olt |itl  itc  itr| ort
//         |             |
//     olc |icl  icc  icr| ort
//         |             |
//     olb |ibl  ibc  ibr| orb
//         +-------------+ 
//          obl  obc  obr
//
type Key struct {
	Hide    bool       // Don't show key/legend if true
	Cols    int        // Number of colums to use. If <0 fill rows before colums
	Border  int        // -1: off, 0: std, 1...:other styles
	Pos     string     // default "" is "itr"
	Entries []KeyEntry // List of entries in the legend
	X, Y    int
}


// KeyEntry encapsulates an antry in the key/legend.
type KeyEntry struct {
	Style  DataStyle
	Symbol int    // Symbol index to use
	Linie  int    // Line Style
	Text   string // Text to display
}

// Margins
var KL_LRBorder int = 1 // before and after whole key
var KL_SLSep int = 2    // space between symbol and test
var KL_ColSep int = 2   // space between columns
var KL_MLSep int = 1    // extra space between rows if multiline text are present

func (key *Key) LayoutKeyTxt() (kb *TextBuf) {
	// TODO(vodo) the following is ugly (and stinks)
	if key.Hide {
		return
	}

	// count real entries in num, see if multilines are present in haveml
	num, haveml := 0, false
	for _, e := range key.Entries {
		if e.Text == "" {
			continue
		}
		num++
		lines := strings.Split(e.Text, "\n", -1)
		if len(lines) > 1 {
			haveml = true
		}
	}
	if num == 0 {
		return
	} // no entries

	rowfirst := false
	cols := key.Cols
	if cols < 0 {
		cols = -cols
		rowfirst = true
	}
	if cols == 0 {
		cols = 1
	}
	if num < cols {
		cols = num
	}
	rows := (num + cols - 1) / cols

	// fmt.Printf("%d entries on %d columns: %d rows\n", num, cols, rows)

	// Arrays with infos
	width := make([][]int, cols)
	for i := 0; i < cols; i++ {
		width[i] = make([]int, rows)
	}
	height := make([][]int, cols)
	for i := 0; i < cols; i++ {
		height[i] = make([]int, rows)
	}
	symbol := make([][]int, cols)
	for i := 0; i < cols; i++ {
		symbol[i] = make([]int, rows)
	}
	text := make([][][]string, cols)
	for i := 0; i < cols; i++ {
		text[i] = make([][]string, rows)
	}

	// fill arrays
	i := 0
	for _, e := range key.Entries {
		if e.Text == "" {
			continue
		}
		var r, c int
		if rowfirst {
			r, c = i/cols, i%cols
		} else {
			c, r = i/rows, i%rows
		}
		lines := strings.Split(e.Text, "\n", -1)
		ml := 0
		for _, t := range lines {
			if len(t) > ml { // TODO(vodo) use utf8.CountRuneInString and honour different chars
				ml = len(t)
			}
		}
		symbol[c][r] = e.Symbol // TODO(vodo) allow line symbols?
		height[c][r] = len(lines)
		width[c][r] = ml
		text[c][r] = lines
		i++
	}
	colwidth := make([]int, cols)
	rowheight := make([]int, rows)
	totalheight, totalwidth := 0, 0
	for c := 0; c < cols; c++ {
		max := 0
		for r := 0; r < rows; r++ {
			if width[c][r] > max {
				max = width[c][r]
			}
		}
		max += 2*KL_LRBorder + 1 + KL_SLSep // formt is " *  Label "
		colwidth[c] = max
		totalwidth += max
	}
	for r := 0; r < rows; r++ {
		max := 0
		for c := 0; c < cols; c++ {
			if height[c][r] > max {
				max = height[c][r]
			}
		}
		rowheight[r] = max
		totalheight += max
	}

	// width and height: + 2 for outer border/box
	w := totalwidth + KL_ColSep*(cols-1) + 2
	h := totalheight + 2
	if haveml {
		h += KL_MLSep * (rows - 1)
	}
	kb = NewTextBuf(w, h)
	if key.Border != -1 {
		kb.Rect(0, 0, w-1, h-1, key.Border+1, ' ')
	}

	// Produce box
	x := 1
	for c := 0; c < cols; c++ {
		y := 1
		for r := 0; r < rows; r++ {
			if width[c][r] == 0 {
				continue
			}
			xx := x + KL_LRBorder
			if symbol[c][r] != -1 {
				kb.Put(xx, y, symbol[c][r])
				xx += 1 + KL_SLSep
			}
			for l, t := range text[c][r] {
				kb.Text(xx, y+l, t, -1)
			}
			y += rowheight[r]
			if haveml {
				y += KL_MLSep
			}
		}
		x += colwidth[c] + KL_ColSep
	}

	return
}


// Place layouts the entries in k in the requested matrix format
func (key Key) Place() (matrix [][]*KeyEntry) {
	// count real entries in num, see if multilines are present in haveml
	num := 0
	for _, e := range key.Entries {
		if e.Text == "" {
			continue
		}
		num++
	}
	if num == 0 {
		return // no entries
	}

	rowfirst := false
	cols := key.Cols
	if cols < 0 {
		cols = -cols
		rowfirst = true
	}
	if cols == 0 {
		cols = 1
	}
	if num < cols {
		cols = num
	}
	rows := (num + cols - 1) / cols

	// Prevent empty last columns in the following case where 5 elements are placed
	// columnsfirst into 4 columns
	//  Col   0    1    2    3
	//       AAA  CCC  EEE
	//       BBB  DDD
	if !rowfirst && rows*(cols-1) >= num {
		cols--
	}

	// Arrays with infos
	matrix = make([][]*KeyEntry, cols)
	for i := 0; i < cols; i++ {
		matrix[i] = make([]*KeyEntry, rows)
	}

	i := 0
	for _, e := range key.Entries {
		if e.Text == "" {
			continue
		}
		var r, c int
		if rowfirst {
			r, c = i/cols, i%cols
		} else {
			c, r = i/rows, i%rows
		}
		matrix[c][r] = &KeyEntry{Text: e.Text, Style: e.Style}
		// fmt.Printf("Place1 (%d,%d) = %d: %s\n", c,r, i, matrix[c][r].Text)
		i++
	}
	return
}


func textviewlen(t string) (length float32) {
	n := 0
	for _, rune := range t {
		if w, ok := CharacterWidth[rune]; ok {
			length += w
		} else {
			length += 23 // save above average
		}
		n++
	}
	length /= averageCharacterWidth
	// fmt.Printf("Length >%s<: %d runes = %.2f  (%d)\n", t, n, length, int(100*length/float32(n)))
	return
}

func textDim(t string) (w float32, h int) {
	lines := strings.Split(t, "\n", -1)
	for _, t := range lines {
		tvl := textviewlen(t)
		if tvl > w {
			w = tvl
		}
	}
	h = len(lines)
	return
}

var (
	KeyColSep      float32 = 2.0
	KeyHorSep      float32 = 1.5
	KeySymbolWidth int     = 30
	KeySymbolSep   int     = 10
	KeyRowSep      float32 = 0.75
	KeyVertSep     float32 = 0.5
)

func (key Key) Layout(bg BasicGraphics, m [][]*KeyEntry) (w, h int, colwidth, rowheight []int) {
	fontwidth, fontheight, _ := bg.FontMetrics(bg.Style("key"))
	cols, rows := len(m), len(m[0])

	// Find total width and height
	totalh := 0
	rowheight = make([]int, rows)
	for r := 0; r < rows; r++ {
		rh := 0
		for c := 0; c < cols; c++ {
			e := m[c][r]
			if e == nil {
				continue
			}
			// fmt.Printf("Layout1 (%d,%d): %s\n", c,r,e.Text)
			_, h := textDim(e.Text)
			if h > rh {
				rh = h
			}
		}
		rowheight[r] = rh
		totalh += rh
	}

	totalw := 0
	colwidth = make([]int, cols)
	// fmt.Printf("Making totalw for %d cols\n", cols)
	for c := 0; c < cols; c++ {
		var rw float32
		for r := 0; r < rows; r++ {
			e := m[c][r]
			if e == nil {
				continue
			}
			// fmt.Printf("Layout2 (%d,%d): %s\n", c,r,e.Text)

			w, _ := textDim(e.Text)
			if w > rw {
				rw = w
			}
		}
		irw := int(rw + 0.75)
		colwidth[c] = irw
		totalw += irw
		// fmt.Printf("Width of col %d: %d.  Total now: %d\n", c, irw, totalw)
	}

	// totalw/h are characters only and still in character-units
	totalw = int(float32(totalw) * fontwidth)                // scale to pixels
	totalw += int(KeyColSep * (float32(cols-1) * fontwidth)) // add space between columns
	totalw += int(2 * KeyHorSep * fontwidth)                 // add space for left/right border
	totalw += (KeySymbolWidth + KeySymbolSep) * cols         // place for symbol and symbol-text sep

	totalh *= fontheight
	totalh += int(KeyRowSep * float32((rows-1)*fontheight)) // add space between rows
	totalh += int(2 * KeyVertSep * float32(fontheight))     // add border at top/bottom

	return totalw, totalh, colwidth, rowheight
}

func GenericKey(bg BasicGraphics, x, y int, key Key) {
	m := key.Place()
	fw, fh, _ := bg.FontMetrics(bg.Style("key"))
	tw, th, cw, rh := key.Layout(bg, m)
	style := bg.Style("key")
	GenericRect(bg, x, y, tw, th, style)
	x += int(KeyHorSep * fw)
	y += int(KeyVertSep*float32(fh)) + fh/2
	for ci, col := range m {
		yy := y

		for ri, e := range col {
			if e == nil || e.Text == "" {
				continue
			}
			s, l, t := e.Style.Symbol, e.Style.LineStyle, e.Text
			// fmt.Printf("Symbol %d=%c, Line=%d: %s\n", s, s, l, t)
			if s == -1 {
				// heading only...
				bg.Text(x, yy, t, "cl", 0, e.Style)
			} else {
				// normal entry
				if l > 0 {
					bg.Line(x, yy, x+KeySymbolWidth, yy, e.Style)
				}
				if s > 0 {
					bg.Symbol(x+KeySymbolWidth/2, yy, s, e.Style)
				}
				bg.Text(x+KeySymbolWidth+KeySymbolSep, yy, t, "cl", 0, e.Style)
			}
			{
				/*
					xx := x + int(fw*float32(cw[ci]))
					bg.Text(xx,y, "|", "cc", 0, e.Style)
					xx += int(KeyColSep*fw)
					bg.Text(xx,y, "|", "cc", 0, e.Style)
					xx += KeySymbolWidth
					bg.Text(xx,y, "|", "cc", 0, e.Style)
					xx += KeySymbolSep
					bg.Text(xx,y, "|", "cc", 0, e.Style)
				*/
			}
			yy += fh*rh[ri] + int(KeyRowSep*float32(fh))
		}

		x += KeySymbolWidth + KeySymbolSep + int(fw*float32(cw[ci])) + int(KeyColSep*fw)
	}
}
