package mqtt_client

import (
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"net/url"
	"time"
)

func main() {
	mqttConfig := mqtt.MQTTConfig{
		Host:     mqtt.MQTTHost,
		Port:     mqtt.MQTTPort,
		User:     url.UserPassword(mqtt.MQTTUsername, mqtt.MQTTPassw),
		Pretopic: mqtt.MQTTPreTopic,
		ClientID: mqtt.MQTTCLientID,
	}
	topics := []string{"vehicle"}
	Publish(topics[0], mqttConfig)
}

func Publish(topic string, mqttConf mqtt.MQTTConfig) {
	client := mqtt.ConnectToBroker("RaspberryPi", mqttConf)
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish(mqttConf.Pretopic+topic, 0, false, t.String())
	}
}
