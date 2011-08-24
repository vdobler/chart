include $(GOROOT)/src/Make.inc


TARG=github.com/vdobler/chart/svgg
GOFILES=\
	svg.go

include $(GOROOT)/src/Make.pkg

format: $(GOFILES)   
	gofmt -w $^