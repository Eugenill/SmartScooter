package utils

import (
	"math"
)

func SoftMax(x []float32) []float64 {
	var max = float64(x[0])
	for _, n := range x {
		max = math.Max(max, float64(n))
	}

	a := make([]float64, len(x))

	var sum float64 = 0
	for i, n := range x {
		a[i] -= math.Exp(float64(n) - max)
		sum += a[i]
	}

	for i, n := range a {
		a[i] = n / sum
	}
	return a
}
