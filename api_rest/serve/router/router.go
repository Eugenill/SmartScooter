package router

import (
	"github.com/Eugenill/SmartScooter/api_rest/serve/handlers"
	"github.com/go-chi/chi"
)

var router *chi.Mux

func SetRouter(name string) *chi.Mux {
	router = chi.NewMux()
	router.Get("/vehicle", handlers.GetVehicles(name))

	return router
}
