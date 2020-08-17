package mqtt_client

import (
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"time"
)

func PublishTimer(topic string, mqttConf mqtt_sub.MQTTConfig) {
	client := mqtt_sub.ConnectToBroker("Timer", mqttConf)
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish(mqttConf.PreTopic+topic, 0, false, t.String())
	}
}
