package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/Eugenill/SmartScooter/api_rest/router"
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
	ctx := &gin.Context{}
	//2. Open DB?
	db.DB = db.OpenDB()

	//3. Initialize MQTT
	client := mqtt_sub.ConnectToBroker(mqttConfig.ClientID, mqttConfig)
	log.Printf("Connected to the MQTTbroker, ClientID:%s, Host:%s, Port%s", mqttConfig.ClientID, mqttConfig.Host, mqttConfig.Port)
	topics := []string{"timer", "detection_example", "RP1_detection"}
	mqtt_sub.ListenToTopics(mqttConfig, topics, client)
	//mqtt_client.PublishTimer("timer", mqttConfig)

	//4. Init server
	engine := router.SetServer(mqttConfig, ctx)
	err := engine.Run("localhost:1234")
	errors.PanicError(err)

	// Shutdown the servers

}
