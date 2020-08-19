package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/endpoints"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/auth"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/Eugenill/SmartScooter/api_rest/server"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
)

var mqttConfig = mqtt_sub.MQTTConfig{
	Host:     mqtt_sub.MQTTHost,
	Port:     mqtt_sub.MQTTPort,
	User:     url.UserPassword(mqtt_sub.MQTTUsername, mqtt_sub.MQTTPassw),
	PreTopic: mqtt_sub.MQTTPreTopic,
	ClientID: mqtt_sub.MQTTCLientID,
}

func main() {
	//1. Create General Context
	//ctx := &gin.Context{}

	//2. Open DB?
	db.DB = db.OpenDB()

	//3. Initialize MQTT
	client := mqtt_sub.ConnectToBroker(mqttConfig.ClientID, mqttConfig)
	log.Printf("Connected to the MQTTbroker, ClientID:%s, Host:%s, Port%s", mqttConfig.ClientID, mqttConfig.Host, mqttConfig.Port)
	topics := []string{"detection_example", "RP1_detection"}
	mqtt_sub.ListenToTopics(mqttConfig, topics, client)

	//4. Init router
	router := server.Router{
		Engine: gin.New(),
		Middlewares: []gin.HandlerFunc{
			gin.Recovery(),
			gin.Logger(),
		},
		MqttConfig: mqttConfig,
	}
	router.AddMiddlewares()

	//Admin Group
	router.AdminRoute = router.Engine.Group("/admin")
	router.AdminRoute.Use(auth.AdminBasicAuth())
	endpoints.AddAdmin(router.AdminRoute, router.MqttConfig)

	//Public Group
	router.PublicRoute = router.Engine.Group("/v1/api")
	router.PublicRoute.Use(auth.PublicBasicAuth())
	endpoints.AddPublic(router.PublicRoute)

	err := router.RunServer()
	errors.PanicError(err)

	// Shutdown the servers

}
