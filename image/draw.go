package main


import (
	"image"
	"image/png"
	"cairo"
	"os"
	"fmt"
)


func main() {
	surface := cairo.NewSurface(cairo.FormatArgB32, 765, 35);
	surface.SelectFontFace("Bitstream Vera Sans Mono", cairo.FontSlantNormal, cairo.FontWeightBold);
	surface.SetFontSize(13.0);
	surface.SetSourceRGB(0.0, 0.0, 0);

	// from U+0021  to U+007E
	surface.MoveTo(0.0, 12);
	text := make([]int, 0, 200)
	for c:='!'; c<='~'; c++ {
		text = append(text, c)
	}
	str := string(text)
	surface.ShowText(str);

	// from U+00A1 to U+00FF
	surface.MoveTo(0.0, 30);
	text = make([]int, 0, 200)
	for c:='¡'; c<='ÿ'; c++ {
		text = append(text, c)
	}
	str = string(text)
	surface.ShowText(str);


	surface.Finish();
	surface.WriteToPNG("font.png");

	ff, err := os.Open("font.png")
	if err!=nil {
		fmt.Printf("Cannot read font.png: %s\n", err.String())
		os.Exit(1)
	}
	
	img, err := png.Decode(ff)
	if err!=nil {
		fmt.Printf("Cannot decode font.png: %s\n", err.String())
		os.Exit(1)
	}

	fg, err := os.Create("font.go")
	if err!=nil {
		fmt.Printf("Cannot create font.go: %s\n", err.String())
		os.Exit(1)
	}


	fmt.Fprintf(fg, "package imgg\n\n")
	fmt.Fprintf(fg, "// Bitstream Vera Sans Mono Bold 13 as 4bit grayscale.\n")
	fmt.Fprintf(fg, "var font map[int][15]uint64 = map[int][15]uint64{\n")
	fmt.Fprintf(fg, "\t%d: [15]uint64{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},\n", ' ')
	for c:='!'; c<='~'; c++ {
		fmt.Fprintf(fg, "\t%d: [15]uint64{\n", c)
		for y:=0; y<15; y++ {
			s := "0x"
			for x:=0; x<8; x++ {
				xx := (c-'!')*8 + x
				var col image.Color
				col = img.At(xx,y)
				_,_,_,a := col.RGBA()
				k := a >> 8
				s += fmt.Sprintf("%02x", k)
			}
			if y<14 {
				fmt.Fprintf(fg, "\t\t%s,\n", s)
			} else {
				fmt.Fprintf(fg, "\t\t%s},\n", s)
			}
		}
	}


	for c:='¡'; c<='ÿ'; c++ {
		fmt.Fprintf(fg, "\t%d: [15]uint64{\n", c)
		for y:=33-15; y<33; y++ {
			s := "0x"
			for x:=0; x<8; x++ {
				xx := (c-'¡')*8 + x
				var col image.Color
				col = img.At(xx,y)
				_,_,_,a := col.RGBA()
				k := a >> 8
				s += fmt.Sprintf("%02x", k)
			}
			if y<32 {
				fmt.Fprintf(fg, "\t\t%s,\n", s)
			} else {
				fmt.Fprintf(fg, "\t\t%s},\n", s)
			}
		}
	}
	fmt.Fprintf(fg,"}\n")
	fg.Close()
}

