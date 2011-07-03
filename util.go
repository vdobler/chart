package chart

import (
	"math"
)

func minimum(data []float64) float64 {
	if data == nil {
		return math.NaN()
	}
	min := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] < min {
			min = data[i]
		}
	}
	return min
}

func maximum(data []float64) float64 {
	if data == nil {
		return math.NaN()
	}
	max := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
		}
	}
	return max
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func clip(x, l, u int) int {
	if x < min(l, u) {
		return l
	}
	if x > max(l, u) {
		return u
	}
	return x
}

func almostEqual(a, b float64) bool {
	rd := math.Fabs((a - b) / (a + b))
	return rd < 1e-5
}
