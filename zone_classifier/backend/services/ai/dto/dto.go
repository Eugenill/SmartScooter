package dto

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}