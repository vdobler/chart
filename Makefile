include $(GOROOT)/src/Make.inc

DRIVERS=\
	svg.go

TARG=github.com/vdobler/chart
GOFILES=\
	chart.go\
	data.go\
	util.go\
	style.go\
	graphics.go\
	stat.go\
	time.go\
	text.go\
	strip.go\
	scatter.go\
	box.go\
	hist.go\
	bar.go\
	cbar.go\
	pie.go\
	$(DRIVERS)

include $(GOROOT)/src/Make.pkg

chart.$(O): $(GOFILES)
	$(GC) -o chart.$(O) $^

samplechart: samplecharts.go chart.$(O)
	$(GC) -I. samplecharts.go
	$(LD) -L. -o samplecharts samplecharts.$(O)

format: $(GOFILES) samplecharts.go  
	gofmt -w $^