include $(GOROOT)/src/Make.inc


TARG=github.com/vdobler/chart/txtg
GOFILES=\
	text.go

include $(GOROOT)/src/Make.pkg

format: $(GOFILES)   
	gofmt -w $^