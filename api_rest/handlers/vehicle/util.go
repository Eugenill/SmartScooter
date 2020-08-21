package vehicle

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
)

//this function context must contain the DB
func CheckVehicle(ctx context.Context, id models.VehicleID) (bool, error) {
	_, err := models.FindVehicle(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

type ReqVehicle struct {
	ID          models.VehicleID   `json:"id"`
	NumberPlate string             `json:"number_plate" `
	HelmetID    models.HelmetID    `json:"helmet_id" `
	IotDevID    models.IotDeviceID `json:"iot_dev_id"`
}
