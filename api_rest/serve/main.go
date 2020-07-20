package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/errors"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt_client"
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

	topics := []string{"vehicle"}

	for _, topic := range topics {
		mqtt.SubscribeToTopic(mqttConfig, topic)
	}
	mqtt_client.Publish("vehicle", mqttConfig)

	//4. Init router
	err := http.ListenAndServe("localhost:1234", router.SetRouter("Ford Mustang"))
	errors.Catch(err)
}
