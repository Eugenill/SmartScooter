package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_client"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/go-chi/chi"
)

func AddEndpoints(router *chi.Mux, mqttConf mqtt_sub.MQTTConfig) {
	router.Get("/publish_detection", mqtt_client.PublishDetection(mqttConf, "detection"))
	router.Get("/vehicle", handlers.GetVehicles("Ford Mustang"))
	router.Post("/create_user", handlers.CreateUser())
}
