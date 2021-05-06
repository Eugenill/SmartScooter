package utils

import (
	"log"
	"sort"
)

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

func ReturnLabels(probabilities []float64, labels []string) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p * 100})
	}
	return higherFirst(resultLabels)
}

// Results implements sort.Interface based on the Probability field.
type Results []LabelResult

func (a Results) Len() int           { return len(a) }
func (a Results) Less(i, j int) bool { return a[i].Probability < a[j].Probability }
func (a Results) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func higherFirst(labelResults []LabelResult) []LabelResult {
	sort.Sort(sort.Reverse(Results(labelResults)))
	log.Print(labelResults)
	return labelResults
}

