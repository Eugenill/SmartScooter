package main

import (
	"bytes"
	"encoding/json"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	"io"
	"log"
	"net/http"
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

const (
	H, W     = 180, 180
	Mean     = float32(45)
	Scale    = float32(1)
	v1_input = "serving_default_rescaling_1_input"
	v2_input = "serving_default_sequential_1_input"
	output   = "StatefulPartitionedCall"
)

var (
	savedModel *tf.SavedModel
	labels     []string
	err        error
)

func main() {
	labels = append(labels, "daisy", "dandelion", "roses", "sunflowers", "tulips")

	savedModel, err = tf.LoadSavedModel("ai/image_classifier/flower_model_v2", []string{"serve"}, nil)
	log.Print("Model Loaded")
	if err != nil {
		log.Panicf("Could not load model files into tensorflow with error: %v", err)
	}

	defer savedModel.Session.Close()

	http.HandleFunc("/prediction", func(w http.ResponseWriter, r *http.Request) {
		predict(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

func makeTensorFromImage(imageBuffer *bytes.Buffer, imageFormat string) (*tf.Tensor, error) {
	tensor, err := tf.NewTensor(imageBuffer.String())
	if err != nil {
		return nil, err
	}
	graph, input, output, err := makeTransformImageGraph(imageFormat)
	if err != nil {
		return nil, err
	}
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func makeTransformImageGraph(imageFormat string) (graph *tf.Graph, input, output tf.Output, err error) {

	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	// Decode PNG or JPEG
	var decode tf.Output
	log.Printf("Image format: %s", imageFormat)
	if imageFormat == "png" {
		decode = op.DecodePng(s, input, op.DecodePngChannels(3))
	} else {
		decode = op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))
	}
	// Div and Sub perform (value-Mean)/Scale for each pixel
	output = op.Div(s,
		op.Sub(s,
			// Resize to 180x180 with bilinear interpolation
			op.ResizeBilinear(s,
				// Create a batch containing a single image
				op.ExpandDims(s,
					// Use decoded pixel values
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func predict(w http.ResponseWriter, r *http.Request) {
	// Read image
	imageFile, header, err := r.FormFile("image")
	// Will contain filename and extension
	imageName := strings.Split(header.Filename, ".")
	if err != nil {
		log.Printf("Could not read image Status: %d", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()
	var imageBuffer bytes.Buffer
	// Copy image data to a buffer
	io.Copy(&imageBuffer, imageFile)

	// Make tensor

	tensor, err := makeTensorFromImage(&imageBuffer, imageName[1])
	if err != nil {
		log.Printf("Invalid image Status: %d", http.StatusBadRequest)
		w.WriteHeader(500)
		return
	}
	log.Print("Tensor from image created")
	//fakeTensor, _ := tf.NewTensor([1][180][180][3]float32{})
	log.Print(tensor)
	s, err := calculate(w, tensor)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error predicting using tf model: %v", err)
		return
	}

	prediction, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling tf output to json: %v", err)
		return
	}
	log.Printf("Prediction: %s", prediction)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(prediction)
}

func calculate(w http.ResponseWriter, tensor *tf.Tensor) (interface{}, error) {
	// Run inference
	result, err := savedModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			savedModel.Graph.Operation(v2_input).Output(0): tensor,
		},
		[]tf.Output{
			savedModel.Graph.Operation(output).Output(0),
		},
		nil)

	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error predicting using tf model: %v", err)
		return nil, err
	}
	return returnLabels(result[0].Value().([][]float32)[0]), err

}

func returnLabels(probabilities []float32) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Return top 5 labels
	return resultLabels
}
