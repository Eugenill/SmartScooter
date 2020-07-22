package router

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt_client"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/log"
	"github.com/go-chi/chi"
)

var router *chi.Mux

func SetRouter(mqttConf mqtt.MQTTConfig) *chi.Mux {
	router = chi.NewMux()
	addEndpoints(router, mqttConf)

	return router
}

func addEndpoints(router *chi.Mux, mqttConf mqtt.MQTTConfig) {
	router.Get("/vehicle", log.AddReqID(handlers.GetVehicles("Ford Mustang")))
	router.Post("/receive_detection_test", log.AddReqID(mqtt_client.PublishDetection(mqttConf, "detection")))
}
