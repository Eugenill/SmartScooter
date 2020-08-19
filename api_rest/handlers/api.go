package handlers

import (
	"github.com/Eugenill/SmartScooter/api_rest/pkg/writters"
	"github.com/gin-gonic/gin"
)

func GetVehicles(name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writters.JsonResponse(ctx, name, 200)
	}
}
