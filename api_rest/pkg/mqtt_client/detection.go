package mqtt_client

import (
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/writters"
	"github.com/gin-gonic/gin"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	"time"
)

func newDetection() string {
	detection := models.Detection{
		TrafficLight: _import00.String{
			String: "red",
			Valid:  true,
		},
		Obstacle:   _import00.String{},
		Location:   models.Point{Latitude: 42.02514, Longitude: 2.055458, Accuracy: 3.5},
		DetectedAt: time.Time{},
	}
	return fmt.Sprintf("%s/%s/%v/%s", detection.TrafficLight.String, detection.Obstacle.String, detection.Location, detection.DetectedAt.String())
}

//handle the req (ping of rasp) and publishes a detection
func PublishDetection(mqttConf mqtt_sub.MQTTConfig, topic string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		detection := newDetection()
		client := mqtt_sub.ConnectToBroker("RaspberryPi", mqttConf)
		client.Publish(mqttConf.PreTopic+topic, 0, false, detection)
		writters.JsonResponse(ctx, detection, 200)
	}
}
