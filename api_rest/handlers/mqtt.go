package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func SaveDetection(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Detection arrived in topic: %s", msg.Topic())
	LogMessage(client, msg)

}

func LogMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}
