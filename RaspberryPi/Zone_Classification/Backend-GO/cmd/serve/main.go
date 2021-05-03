package main

import (
	"github.com/Eugenill/SmartScooter/RaspBerryPi/Zone_Classifier/Backend-GO/endpoints"
	"github.com/Eugenill/SmartScooter/RaspBerryPi/Zone_Classifier/Backend-GO/pkg/errors"
	"github.com/Eugenill/SmartScooter/RaspBerryPi/Zone_Classifier/Backend-GO/server"
	"github.com/gin-gonic/gin"
)

func main() {
	//Init router
	router := server.Router{
		Engine: gin.New(),
		Middlewares: []gin.HandlerFunc{
			gin.Recovery(),
			gin.Logger(),
		},
	}
	router.AddMiddlewares()

	//Public Group
	router.PublicRoute = router.Engine.Group("/v1/api")
	endpoints.AddPublic(router.PublicRoute)

	err := router.RunServer()
	errors.PanicError(err)
}
