package endpoints

import (
	"github.com/Eugenill/SmartScooter/RaspBerryPi/Zone_Classifier/Backend-GO/handlers/image_classifier"
	"github.com/gin-gonic/gin"
)

func AddPublic(group *gin.RouterGroup) {
	group.POST("/prediction", image_classifier.MakePrediction())
}
