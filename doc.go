/*
Package chart implements common chart/plot types.

The following chart types are available:

	StripChart		Visualize set of numeric values
	ScatterChart   		Plot (x,y) data (with optional error bars) and/or functions
	HistChart		Produce histograms from data
	BarChart		Show (x,y) data as bars
	CategoryBarChart	Bar chart of categorical (non-numeric) data
	BoxChart            	Box charts to visualize distributions
	PieChart            	Pie and Ring charts

Chart tries to provides useful defaults and produce nice charts without sacrificing accuracy.
The generated charts look good and are higly customizable but will not match handmade
photoshop charts done by marketing.

Creating charts consists of the following steps:
	1. Create chart object
	2. Configure chart, axis, etc.
	3. Add data
	4. Render chart to graphic output
You may change the configuration at any step or render to different outputs.



*/
package chart
