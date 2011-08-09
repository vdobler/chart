include $(GOROOT)/src/Make.inc


TARG=github.com/vdobler/chart/txt
GOFILES=\
	text.go

include $(GOROOT)/src/Make.pkg

format: $(GOFILES)   
	gofmt -w $^