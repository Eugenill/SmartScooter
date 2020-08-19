package vehicle

import (
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/contxt"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddHelmet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vID := contxt.RequestHeader(ctx, "vehicleID")
		vehID, err := models.VehicleIDFromString(vID)
		if err != nil {
			_, ginErr := errors.New(ctx, "vehicleId is not valid", gin.ErrorTypePrivate)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
		ok, err := CheckVehicle(ctx, vehID)
		if !ok {
			_, ginErr := errors.New(ctx, "vehicle not exist", gin.ErrorTypePrivate)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
	}
}
