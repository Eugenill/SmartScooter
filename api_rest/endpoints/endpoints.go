package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/log"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_client"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/go-chi/chi"
)

func AddEndpoints(router *chi.Mux, mqttConf mqtt_sub.MQTTConfig) {
	router.Get("/vehicle", log.AddReqID(handlers.GetVehicles("Ford Mustang")))
	router.Get("/publish_detection", log.AddReqID(mqtt_client.PublishDetection(mqttConf, "detection")))
}
