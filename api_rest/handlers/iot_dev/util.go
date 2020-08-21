package iot_dev

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
)

//this function context must contain the DB
func CheckIotDev(ctx context.Context, id models.IotDeviceID) (bool, error) {
	_, err := models.FindIotDevice(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

type CreateIotDev struct {
	Name string `json:"name"`
}
