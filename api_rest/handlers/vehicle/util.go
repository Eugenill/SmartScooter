package vehicle

import (
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/gin-gonic/gin"
)

func CheckVehicle(ctx *gin.Context, id models.VehicleID) (bool, error) {
	_, err := models.FindVehicle(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
