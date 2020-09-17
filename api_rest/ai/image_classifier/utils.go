package image_classifier

import "math"

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

func returnLabels(probabilities []float64) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p * 100})
	}
	return resultLabels
}
