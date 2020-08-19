package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/gin-gonic/gin"
)

func AddPublic(group *gin.RouterGroup) {
	group.GET("/vehicle", handlers.GetVehicles("Ford Mustang"))
}
