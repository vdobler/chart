include $(GOROOT)/src/Make.inc

TARG=github.com/vdobler/chart
GOFILES=\
	chart.go\
	text.go\
	strip.go\
	scatter.go

include $(GOROOT)/src/Make.pkg

chart.$(O): chart.go text.go strip.go scatter.go
	$(GC) -o chart.$(O) $^

samplechart: samplecharts.go chart.$(O)
	$(GC) -I. samplecharts.go
	$(LD) -L. -o samplecharts samplecharts.$(O)

format: $(GOFILES) samplecharts.go  
	gofmt -w $^