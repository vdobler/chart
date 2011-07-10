package chart

import (
// "fmt"
)

var Edge = [][4]int{{'+', '+', '+', '+'}, {'.', '.', '\'', '\''}, {'/', '\\', '\\', '/'}}

var Symbol = []int{'*', '+', 'o', '#', '='}

type TextBuf struct {
	Buf  []int
	W, H int
}

// Data is from 0 to w-1. Pos w is nl.
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

func (tb *TextBuf) Put(x, y, c int) {
	if x < 0 || y < 0 || x >= tb.W || y >= tb.H {
		// fmt.Printf("Falsch: %d, %d  '%c' \n", x, y, c)
		x, y = 0, 0
	}
	i := (tb.W+1)*y + x
	tb.Buf[i] = c
}

// buf[x+y*s] is pos x,y
func (tb *TextBuf) Rect(x, y, w, h int, style int, fill int) {
	style = style % 3

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


func StrLen(s string) (n int) {
	for _, _ = range s {
		n++
	}
	return
}

// align: -1: left, 0: centered, 1: right, 2 |top, 3 |center, 4 |bot
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


func (tb *TextBuf) Paste(x, y int, buf *TextBuf) {
	s := buf.W + 1
	for i := 0; i < buf.W; i++ {
		for j := 0; j < buf.H; j++ {
			tb.Put(x+i, y+j, buf.Buf[i+s*j])
		}
	}
}


func (tb *TextBuf) String() string {
	return string(tb.Buf)
}
