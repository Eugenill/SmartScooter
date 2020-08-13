package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func SaveDetection(client mqtt.Client, msg mqtt.Message) {
	log.Println("Detection arrived")
	log.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))

}
func LogDetection(client mqtt.Client, msg mqtt.Message) {
	log.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}
