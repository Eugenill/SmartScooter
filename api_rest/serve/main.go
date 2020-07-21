package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/errors"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"github.com/Eugenill/SmartScooter/api_rest/router"
	"net/http"
	"net/url"
)

func main() {
	//1. Create Context

	//2. Open DB

	//3. Initialize MQTT
	mqttConfig := mqtt.MQTTConfig{
		Host:     mqtt.MQTTHost,
		Port:     mqtt.MQTTPort,
		User:     url.UserPassword(mqtt.MQTTUsername, mqtt.MQTTPassw),
		Pretopic: mqtt.MQTTPreTopic,
		ClientID: mqtt.MQTTCLientID,
	}

	topics := []string{"timer", "detection"}
	go mqtt.ListenToTopic(mqttConfig, topics[0])
	go mqtt.ListenToTopic(mqttConfig, topics[1])
	//mqtt_client.PublishTimer("timer", mqttConfig)

	//4. Init router
	err := http.ListenAndServe("localhost:1234", router.SetRouter(mqttConfig))
	errors.Catch(err)
}
