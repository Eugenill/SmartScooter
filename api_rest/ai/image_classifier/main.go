package image_classifier

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
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
	Probability float64 `json:"probability"`
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

func init() {
	labels = append(labels, "daisy", "dandelion", "roses", "sunflowers", "tulips")

	savedModel, err = tf.LoadSavedModel("ai/image_classifier/flower_model_v2", []string{"serve"}, nil)
	log.Print("Model Loaded")
	if err != nil {
		log.Panicf("Could not load model files into tensorflow with error: %v", err)
	}
}

func MakePrediction() gin.HandlerFunc {
	return func(context *gin.Context) {
		r := context.Request
		w := context.Writer
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
	probabilities := SoftMax(result[0].Value().([][]float32)[0])
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error predicting using tf model: %v", err)
		return nil, err
	}
	return returnLabels(probabilities), err

}
