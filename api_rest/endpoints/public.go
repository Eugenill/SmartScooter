package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/ai/image_classifier"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/ride"
	"github.com/gin-gonic/gin"
)

func AddPublic(group *gin.RouterGroup) {
	group.GET("/create_ride", ride.CreateRide())
	group.GET("/finish", ride.FinishRide())
	group.POST("/prediction", image_classifier.MakePrediction())
}
