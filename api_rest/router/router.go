package router

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt"
	"github.com/Eugenill/SmartScooter/api_rest/mqtt_client"
	"github.com/go-chi/chi"
)

var router *chi.Mux

func SetRouter(mqttConf mqtt.MQTTConfig) *chi.Mux {
	router = chi.NewMux()
	name := "Ford Mustang"
	router.Get("/vehicle", handlers.GetVehicles(name))
	router.Post("/receive_detection_test", mqtt_client.PublishDetection(mqttConf, "detection"))

	return router
}
