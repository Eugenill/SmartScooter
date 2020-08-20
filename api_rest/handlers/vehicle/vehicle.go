package vehicle

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/helmet"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"net/http"
)

func AddVehicle() gin.HandlerFunc {
	meta := map[string]string{"function": "AddVehicle"}
	return func(ctxGin *gin.Context) {
		vehicle := &CreateVehicle{}
		r := ctxGin.Request
		ginErr := &gin.Error{}
		var err error
		ctx2 := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&vehicle, r); err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			vehID, err := models.VehicleIDFromString(vehicle.ID)
			if err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			helmetID, err := models.HelmetIDFromString(vehicle.HelmetID)
			if err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			ok, err := helmet.CheckHelmet(ctx2, helmetID)
			if !ok {
				err, ginErr = errors.New(ctxGin, "helmet does not exist", gin.ErrorTypePrivate, meta)
				return err
			}
			vehicle := &models.Vehicle{
				ID:            vehID,
				CurrentRideID: models.NullRideID{},
				LastRideID:    models.NullRideID{},
				CurrentUserID: models.NullUserID{},
				LastUserID:    models.NullUserID{},
				NumberPlate:   vehicle.NumberPlate,
				HelmetID:      helmetID,
			}
			err = vehicle.Insert(ctx2)
			if err != nil {
				_, ginErr = errors.New(ctxGin, "vehicle not inserted", gin.ErrorTypePrivate, meta)
				return err
			}
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginErr, http.StatusBadRequest)
		}
	}
}
