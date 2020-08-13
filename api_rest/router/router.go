package router

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/endpoints"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"github.com/go-chi/chi"
)

var router *chi.Mux

func SetRouter(mqttConf mqtt_sub.MQTTConfig, ctx context.Context) *chi.Mux {
	router = chi.NewMux()
	rest.AddMiddlewares(router)
	endpoints.AddEndpoints(router, mqttConf)

	return router
}
