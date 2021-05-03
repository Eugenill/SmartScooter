package image_classifier

import (
	"bytes"
	"encoding/json"
	"github.com/Eugenill/SmartScooter/RaspBerryPi/Zone_Classifier/Backend-GO/ai/image_classifier/zone_class"
	"github.com/Eugenill/SmartScooter/RaspBerryPi/Zone_Classifier/Backend-GO/pkg/utils"
	"github.com/gin-gonic/gin"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"io"
	"log"
	"net/http"
	"strings"
)

// https://www.tensorflow.org/install/lang_c to build Tensorflow C

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

var (
	savedModel *tf.SavedModel
	labels     []string
	err        error
)

func init() {
	labels = append(labels, "carril_bici", "acera")

	savedModel, err = tf.LoadSavedModel("ai/image_classifier/zone_class", []string{"serve"}, nil)
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
		log.Print(err)
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

		tensor, err := utils.MakeTensorFromImage(&imageBuffer, imageName[1])
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
			savedModel.Graph.Operation(zone_class.Input).Output(0): tensor,
		},
		[]tf.Output{
			savedModel.Graph.Operation(zone_class.Output).Output(0),
		},
		nil)
	log.Print(result[0].Value().([][]float32)[0])
	probabilities := utils.SoftMax(result[0].Value().([][]float32)[0])
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error predicting using tf model: %v", err)
		return nil, err
	}
	return utils.ReturnLabels(probabilities, labels), err

}
