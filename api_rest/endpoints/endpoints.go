package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_client"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/gin-gonic/gin"
)

func AddEndpoints(engine *gin.Engine, mqttConf mqtt_sub.MQTTConfig) {
	engine.GET("/publish_detection", mqtt_client.PublishDetection(mqttConf, "detection_example"))
	engine.GET("/vehicle", handlers.GetVehicles("Ford Mustang"))
	engine.POST("/create_user", handlers.CreateUser())
}
