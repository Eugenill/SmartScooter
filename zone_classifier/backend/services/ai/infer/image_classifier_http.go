package infer

import (
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/ai_models/image_classifier/zone_class"
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/services/ai/utils"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"log"
	"net/http"
)

// https://www.tensorflow.org/install/lang_c to build Tensorflow C
type ImageClassifierModel struct {
	SavedModel *tf.SavedModel
	Labels     []string
}

func (i *ImageClassifierModel)Calculate(w http.ResponseWriter, tensor *tf.Tensor) (interface{}, error) {
	// Run inference
	result, err := i.SavedModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			i.SavedModel.Graph.Operation(zone_class.Input).Output(0): tensor,
		},
		[]tf.Output{
			i.SavedModel.Graph.Operation(zone_class.Output).Output(0),
		},
		nil)
	probabilities := utils.SoftMax(result[0].Value().([][]float32)[0])
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error predicting using infer model: %v", err)
		return nil, err
	}
	return utils.ReturnLabels(probabilities, i.Labels), err

}
