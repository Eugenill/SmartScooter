package router

import (
	"github.com/Eugenill/SmartScooter/api_rest/endpoints"
	"github.com/Eugenill/SmartScooter/api_rest/middleware"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/gin-gonic/gin"
)

func SetServer(mqttConf mqtt_sub.MQTTConfig) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Logger(), gin.Recovery())
	middleware.AddMiddlewares(engine)
	endpoints.AddEndpoints(engine, mqttConf)
	return engine
}
