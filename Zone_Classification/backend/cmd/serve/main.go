package main

import (
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/pkg/errors"
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/pkg/kernel"
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/services/ai"
	"github.com/gin-gonic/gin"
)

func main() {
	//Init kernel
	krnl := kernel.New(gin.New(), []gin.HandlerFunc{gin.Recovery(), gin.Logger()})

	//Add AI Service with its RouterGroup
	krnl.AddService(ai.New(krnl.Engine.Group("/ai")))

	//Init Kernel
	err := krnl.Init()
	errors.PanicError(err)
}
