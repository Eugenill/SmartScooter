package router

import (
	"github.com/Eugenill/SmartScooter/api_rest/endpoints"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/go-chi/chi"
)

var router *chi.Mux

func SetRouter(mqttConf mqtt_sub.MQTTConfig) *chi.Mux {
	router = chi.NewMux()
	endpoints.AddEndpoints(router, mqttConf)

	return router
}
