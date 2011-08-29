include $(GOROOT)/src/Make.inc


TARG=github.com/vdobler/chart/imgg
GOFILES=\
	image.go\
	font.go

include $(GOROOT)/src/Make.pkg

font.go: draw
	./draw

draw: draw.go
	$(GC) draw.go
	$(LD) -o draw draw.$(O)


format: $(GOFILES)   
	gofmt -w $^