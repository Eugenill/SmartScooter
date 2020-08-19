package ride

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"net/http"
	"time"

	"github.com/Eugenill/SmartScooter/api_rest/handlers/auth"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/vehicle"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/contxt"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

func CreateRide() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ginErr *gin.Error
		vID := contxt.RequestHeader(ctx, "vehicleID")
		vehID, err := models.VehicleIDFromString(vID)
		if err != nil {
			_, ginErr = errors.New(ctx, "vehicleID not valid", gin.ErrorTypePrivate)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
		ok, err := vehicle.CheckVehicle(ctx, vehID)
		if !ok {
			_, ginErr = errors.New(ctx, "vehicle not exist", gin.ErrorTypePrivate)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)

		}
		user, err := auth.FromContext(ctx)
		if err != nil {
			_, ginErr = errors.New(ctx, "user not found in context", gin.ErrorTypePrivate)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)

		}

		//Call Ride Reby endpoint
		_, err = http.NewRequest("POST", "https://api.reby.co/v2/research/ride", nil)
		if err != nil {
			_, ginErr = errors.New(ctx, "ride to the scooter failed", gin.ErrorTypePrivate)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}

		//Adding ride and updating vehicle
		ctx2 := db.GinToContextWithDB(ctx)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			now := time.Now()
			ride := &models.Ride{
				ID:        models.NewRideID(),
				VehicleID: vehID,
				UserID:    user.ID,
				StartedAt: now,
			}
			err = ride.Insert(ctx2)
			if err != nil {
				_, ginErr = errors.New(ctx, "ride not inserted", gin.ErrorTypePrivate)
				return err
			}
			veh, err := models.FindVehicle(ctx2, vehID)
			if err != nil {
				_, ginErr = errors.New(ctx, "vehicle not found in db", gin.ErrorTypePrivate, vehID)
				return err
			}

			veh.CurrentRideID = ride.ID
			veh.CurrentUserID = user.ID
			if err = veh.Update(ctx2); err != nil {
				_, ginErr = errors.New(ctx, "vehicle not updated", gin.ErrorTypePrivate, vehID)
				return err
			}
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
	}
}
