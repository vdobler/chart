package chart

import (
// "fmt"
)

// Different edge styles for boxes
var Edge = [][4]int{{'+', '+', '+', '+'}, {'.', '.', '\'', '\''}, {'/', '\\', '\\', '/'}}

// Different symbols
var Symbol = []int{'o', '=', '#', '%', '+', 'X', '*', '@', '$', 'H', 'A', '0', 'V', '.'}

var CharacterWidth = map[int]float64{'a': 16.8, 'b': 17.0, 'c': 15.2, 'd': 16.8, 'e': 16.8, 'f': 8.5, 'g': 17.0,
	'h': 16.8, 'i': 5.9, 'j': 5.9, 'k': 16.8, 'l': 6.9, 'm': 25.5, 'n': 16.8, 'o': 16.8, 'p': 17.0, 'q': 17.0,
	'r': 10.2, 's': 15.2, 't': 8.4, 'u': 16.8, 'v': 15.4, 'w': 22.2, 'x': 15.2, 'y': 15.2, 'z': 15.2,
	'A': 20.2, 'B': 20.2, 'C': 22.2, 'D': 22.2, 'E': 20.2, 'F': 18.6, 'G': 23.5, 'H': 22.0, 'I': 8.2, 'J': 15.2,
	'K': 20.2, 'L': 16.8, 'M': 25.5, 'N': 22.0, 'O': 23.5, 'P': 20.2, 'Q': 23.5, 'R': 21.1, 'S': 20.2, 'T': 18.5,
	'U': 22.0, 'V': 20.2, 'W': 29.0, 'X': 20.2, 'Y': 20.2, 'Z': 18.8,
	'1': 16.8, '2': 16.8, '3': 16.8, '4': 16.8, '5': 16.8, '6': 16.8, '7': 16.8, '8': 16.8, '9': 16.8, '0': 16.8,
	'.': 8.2, ',': 8.2, ':': 8.2, ';': 8.2, '+': 17.9, '"': 11.0, '*': 11.8, '%': 27.0, '&': 20.2, '/': 8.4,
	'(': 10.2, ')': 10.2, '=': 18.0, '?': 16.8, '!': 8.5, '[': 8.2, ']': 8.2, '{': 10.2, '}': 10.2, '$': 16.8,
	'<': 18.0, '>': 18.0, '§': 16.8, '°': 12.2, '^': 14.2, '~': 18.0,
}


// A Text Buffer
type TextBuf struct {
	Buf  []int // the internal buffer.  Point (x,y) is mapped to x + y*(W+1)
	W, H int   // Width and Height
}

// Set up a new TextBuf with width w and height h.
func NewTextBuf(w, h int) (tb *TextBuf) {
	tb = new(TextBuf)
	tb.W, tb.H = w, h
	tb.Buf = make([]int, (w+1)*h)
	for i, _ := range tb.Buf {
		tb.Buf[i] = ' '
	}
	for i := 0; i < h; i++ {
		tb.Buf[i*(w+1)+w] = '\n'
	}
	// tb.Buf[0], tb.Buf[(w+1)*h-1] = 'X', 'X'
	return
}


// Put character c at (x,y)
func (tb *TextBuf) Put(x, y, c int) {
	if x < 0 || y < 0 || x >= tb.W || y >= tb.H {
		return
		// fmt.Printf("Ooooops Put(): %d, %d  '%c' \n", x, y, c)
		x, y = 0, 0
	}
	i := (tb.W+1)*y + x
	tb.Buf[i] = c
}

// Draw rectangle of width w and height h from corner at (x,y).
// Use one of the corner style defined in Edge. 
// Interior is filled with charater fill iff != 0.
func (tb *TextBuf) Rect(x, y, w, h int, style int, fill int) {
	style = style % len(Edge)

	if h < 0 {
		h = -h
		y -= h
	}
	if w < 0 {
		w = -w
		x -= w
	}

	tb.Put(x, y, Edge[style][0])
	tb.Put(x+w, y, Edge[style][1])
	tb.Put(x, y+h, Edge[style][2])
	tb.Put(x+w, y+h, Edge[style][3])
	for i := 1; i < w; i++ {
		tb.Put(x+i, y, '-')
		tb.Put(x+i, y+h, '-')
	}
	for i := 1; i < h; i++ {
		tb.Put(x+w, y+i, '|')
		tb.Put(x, y+i, '|')
		if fill > 0 {
			for j := x + 1; j < x+w; j++ {
				tb.Put(j, y+i, fill)
			}
		}
	}
}

func (tb *TextBuf) Block(x, y, w, h int, fill int) {
	if h < 0 {
		h = -h
		y -= h
	}
	if w < 0 {
		w = -w
		x -= w
	}
	for i := 0; i < w; i++ {
		for j := 0; j <= h; j++ {
			tb.Put(x+i, y+j, fill)
		}
	}
}

// Return real character len of s (rune count).
func StrLen(s string) (n int) {
	for _, _ = range s {
		n++
	}
	return
}

// Print text txt at (x,y). Horizontal display for align in [-1,1],
// vasrtical alignment for align in [2,4]
// align: -1: left; 0: centered; 1: right; 2: top, 3: center, 4: bottom
func (tb *TextBuf) Text(x, y int, txt string, align int) {
	if align <= 1 {
		switch align {
		case 0:
			x -= StrLen(txt) / 2
		case 1:
			x -= StrLen(txt)
		}
		i := 0
		for _, rune := range txt {
			tb.Put(x+i, y, rune)
			i++
		}
	} else {
		switch align {
		case 3:
			y -= StrLen(txt) / 2
		case 2:
			x -= StrLen(txt)
		}
		i := 0
		for _, rune := range txt {
			tb.Put(x, y+i, rune)
			i++
		}
	}
}


// Paste buf at (x,y)
func (tb *TextBuf) Paste(x, y int, buf *TextBuf) {
	s := buf.W + 1
	for i := 0; i < buf.W; i++ {
		for j := 0; j < buf.H; j++ {
			tb.Put(x+i, y+j, buf.Buf[i+s*j])
		}
	}
}

func (tb *TextBuf) Line(x0, y0, x1, y1 int, symbol int) {
	// handle trivial cases first
	if x0 == x1 {
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		for ; y0 <= y1; y0++ {
			tb.Put(x0, y0, symbol)
		}
		return
	}
	if y0 == y1 {
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		for ; x0 <= x1; x0++ {
			tb.Put(x0, y0, symbol)
		}
		return
	}
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := sign(x1-x0), sign(y1-y0)
	err, e2 := dx+dy, 0
	for {
		tb.Put(x0, y0, symbol)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 = 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}

	}
}


// Convert to plain (utf8) string.
func (tb *TextBuf) String() string {
	return string(tb.Buf)
}
