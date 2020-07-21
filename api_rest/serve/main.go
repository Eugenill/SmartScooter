package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/errors"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"github.com/Eugenill/SmartScooter/api_rest/router"
	"log"
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
	client := mqtt.ConnectToBroker(mqttConfig.ClientID, mqttConfig)
	log.Printf("Connected to the MQTTbroker, Host:%s, Port%s", mqttConfig.Host, mqttConfig.Port)
	topics := []string{"timer", "detection"}
	for _, topic := range topics {
		go mqtt.ListenToTopic(mqttConfig, topic, client)
		//mqtt_client.PublishTimer("timer", mqttConfig)
	}
	//4. Init router
	err := http.ListenAndServe("localhost:1234", router.SetRouter(mqttConfig))
	errors.Catch(err)
}
