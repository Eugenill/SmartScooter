package ai

import (
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/pkg/kernel/http/api"
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/services/ai/app"
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/services/ai/infer"
	"github.com/gin-gonic/gin"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"log"
	"net/http"
)

type Service struct {
	RouterGroup     *gin.RouterGroup
	ImageClassifier *infer.ImageClassifierModel
}

func New(routerGroup *gin.RouterGroup) *Service {
	return &Service{
		RouterGroup:     routerGroup,
	}
}

func (s *Service) RegisterHTTPHandlers() []api.HTTPHandlerDefinition {
	return []api.HTTPHandlerDefinition{
		{Method: http.MethodPost, RouterGroup: s.RouterGroup, Endpoint: "/prediction", Handler: app.MakePrediction(s.ImageClassifier)},
	}
}

func (s *Service) Init() {
	labels := append([]string{}, "carril_bici", "acera")
	savedModel, err := tf.LoadSavedModel("ai_models/image_classifier/zone_class", []string{"serve"}, nil)
	log.Print("Model Loaded")
	if err != nil {
		log.Panicf("Could not load model files into tensorflow with error: %v", err)
	}
	s.ImageClassifier = &infer.ImageClassifierModel{
		SavedModel: savedModel,
		Labels: labels,
	}
}
