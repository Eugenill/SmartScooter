package image_classifier

import (
	"bytes"
	"github.com/gin-gonic/gin"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"io"
	"log"
	"net/http"

	"sort"
	"strings"
)

type ClassifyResult struct {
	Filename string        `json:"filename"`
	Labels   []LabelResult `json:"labels"`
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

var (
	graphModel *tf.Graph
	savedModel *tf.SavedModel
	labels     []string
)

func init() {
	if err := loadModel(); err != nil {
		log.Fatal(err)
		return
	}
}

func loadModel() error {
	// Load inception flowers_model
	var err error
	savedModel, err = tf.LoadSavedModel("ai/image_classifier/flowers_model", []string{"serve"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load labels
	labels = append(labels, "daisy", "dandelion", "roses", "sunflowers", "tulips")
	/*
		labelsFile, err := os.Open("/flowers_model/imagenet_comp_graph_label_strings.txt")

		if err != nil {
			return err
		}
		defer labelsFile.Close()
		scanner := bufio.NewScanner(labelsFile)
		// Labels are separated by newlines
		for scanner.Scan() {
			labels = append(labels, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	*/
	return nil
}

func RecognizeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		w := ctx.Writer
		// Read image
		imageFile, header, err := r.FormFile("image")
		// Will contain filename and extension
		imageName := strings.Split(header.Filename, ".")
		if err != nil {
			responseError(w, "Could not read image", http.StatusBadRequest)
			return
		}
		defer imageFile.Close()
		var imageBuffer bytes.Buffer
		// Copy image data to a buffer
		io.Copy(&imageBuffer, imageFile)

		// ...
		// Make tensor
		tensor, err := makeTensorFromImage(&imageBuffer, imageName[:1][0])
		if err != nil {
			responseError(w, "Invalid image", http.StatusBadRequest)
			return
		}

		// Run inference
		output, err := savedModel.Session.Run(
			map[tf.Output]*tf.Tensor{
				graphModel.Operation("input").Output(0): tensor,
			},
			[]tf.Output{
				graphModel.Operation("output").Output(0),
			},
			nil)
		if err != nil {
			responseError(w, "Could not run inference", http.StatusInternalServerError)
			return
		}

		// Return best labels
		responseJSON(w, ClassifyResult{
			Filename: header.Filename,
			Labels:   findBestLabels(output[0].Value().([][]float32)[0]),
		})
	}
}

type ByProbability []LabelResult

func (a ByProbability) Len() int           { return len(a) }
func (a ByProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }

func findBestLabels(probabilities []float32) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Sort by probability
	sort.Sort(ByProbability(resultLabels))
	// Return top 5 labels
	return resultLabels[:5]
}
