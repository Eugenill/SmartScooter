package mqtt_client

import (
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/errors"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	_import01 "github.com/sqlbunny/geo"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	"log"
	"net/http"
	"time"
)

func PublishTimer(topic string, mqttConf mqtt.MQTTConfig) {
	client := mqtt.ConnectToBroker("RaspberryPi", mqttConf)
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish(mqttConf.Pretopic+topic, 0, false, t.String())
	}
}

func newDetection() string {
	detection := models.Detection{
		TrafficLight: _import00.String{
			String: "red",
			Valid:  true,
		},
		Obstacle:   _import00.String{},
		Location:   _import01.Point{},
		DetectedAt: time.Time{},
	}
	return fmt.Sprintf("%s,%s,%v,%s", detection.TrafficLight.String, detection.Obstacle.String, detection.Location, detection.DetectedAt.String())
}

//handle the req (ping of rasp) and publishes a detection
func PublishDetection(mqttConf mqtt.MQTTConfig, topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		detection := newDetection()
		client := mqtt.ConnectToBroker("RaspberryPi", mqttConf)
		client.Publish(mqttConf.Pretopic+topic, 0, false, detection)
		_, err := w.Write([]byte(detection))
		errors.Catch(err)
	}
}
