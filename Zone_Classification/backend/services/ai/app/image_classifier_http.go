package app

import (
	"bytes"
	"encoding/json"
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/services/ai/infer"
	utils2 "github.com/Eugenill/SmartScooter/Zone_Classifier/backend/services/ai/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

// https://www.tensorflow.org/install/lang_c to build Tensorflow C

func MakePrediction(model *infer.ImageClassifierModel) gin.HandlerFunc {
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
		tensor, err := utils2.MakeTensorFromImage(&imageBuffer, imageName[1])
		if err != nil {
			log.Printf("Invalid image Status: %d", http.StatusBadRequest)
			w.WriteHeader(500)
			return
		}
		//fakeTensor, _ := infer.NewTensor([1][180][180][3]float32{})
		s, err := model.Calculate(w, tensor)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("Error predicting using infer model: %v", err)
			return
		}

		prediction, err := json.Marshal(s)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("Error marshalling infer output to json: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(prediction)
	}
}

