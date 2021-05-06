package endpoints

import (
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/handlers/image_classifier"
	"github.com/gin-gonic/gin"
)

func AddPublic(group *gin.RouterGroup) {
	group.POST("/prediction", image_classifier.MakePrediction())
}
