package ride

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/reby"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	"net/http"
	"time"

	"github.com/Eugenill/SmartScooter/api_rest/handlers/auth"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/contxt"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

func CreateRide() gin.HandlerFunc {
	meta := map[string]string{"app": "CreateRide"}
	return func(ctx *gin.Context) {
		var ginErr *gin.Error
		vID := contxt.RequestHeader(ctx, "vehicleID")
		vehID, err := models.VehicleIDFromString(vID)
		if err != nil {
			_, ginErr = errors.New(ctx, "vehicleID not valid", gin.ErrorTypePrivate, meta)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}

		user, err := auth.FromContext(ctx)
		if err != nil {
			_, ginErr = errors.New(ctx, "user not found in context", gin.ErrorTypePrivate, meta)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)

		}

		//Call Ride Reby endpoint
		req, err := http.NewRequest("POST", "https://api.reby.co/v2/research/ride", nil)
		if err != nil {
			_, ginErr = errors.New(ctx, "ride call request creation failed", gin.ErrorTypePrivate, meta)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		} else {
			req.Header = reby.SetHeaders(vID, reby.BearerETSEIB)
		}
		client := http.DefaultClient
		_, err = client.Do(req)
		now := time.Now()
		if err != nil {
			_, ginErr = errors.New(ctx, "ride call to the scooter failed", gin.ErrorTypePrivate, meta)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}

		//Adding ride and updating vehicle
		ctx2 := db.GinToContextWithDB(ctx)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			ride := &models.Ride{
				ID:        models.NewRideID(),
				VehicleID: vehID,
				UserID:    user.ID,
				StartedAt: now,
				PathID:    models.NewPathID(),
			}
			err = ride.Insert(ctx2)
			if err != nil {
				_, ginErr = errors.New(ctx, "ride not inserted", gin.ErrorTypePrivate, meta)
				return err
			}
			veh, err := models.FindVehicle(ctx2, vehID)

			if err != nil {
				_, ginErr = errors.New(ctx, "vehicle not found in db", gin.ErrorTypePrivate, vehID, meta)
				return err
			}

			veh.CurrentRideID = models.NullRideIDFrom(ride.ID)
			veh.CurrentUserID = models.NullUserIDFrom(user.ID)
			if err = veh.Update(ctx2); err != nil {
				_, ginErr = errors.New(ctx, "vehicle not updated", gin.ErrorTypePrivate, vehID, meta)
				return err
			}
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
	}
}

func FinishRide() gin.HandlerFunc {
	meta := map[string]string{"app": "CreateRide"}
	return func(ctx *gin.Context) {
		var ginErr *gin.Error
		rID := contxt.RequestHeader(ctx, "rideID")
		rideID, err := models.RideIDFromString(rID)
		if err != nil {
			_, ginErr = errors.New(ctx, "rideID not valid", gin.ErrorTypePrivate, meta)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
		r, err := models.FindRide(ctx, rideID)
		if err != nil {
			_, ginErr = errors.New(ctx, "ride not found", gin.ErrorTypePrivate, meta)
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}

		//Adding ride and updating vehicle
		ctx2 := db.GinToContextWithDB(ctx)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			//TODO: Add distance and path
			r.PathID = models.PathID{}
			r.Distance = 0

			r.FinishedAt = _import00.Time{
				Time:  time.Now(),
				Valid: true,
			}
			r.Duration = int32(r.FinishedAt.Time.Sub(r.StartedAt))
			err = r.Update(ctx2)
			if err != nil {
				err, ginErr = errors.New(ctx, "ride not updated", gin.ErrorTypePrivate, meta)
				return err
			}
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
	}
}
