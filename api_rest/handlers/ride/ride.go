package ride

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/reby"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/writters"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	"net/http"
	"time"

	"github.com/Eugenill/SmartScooter/api_rest/handlers/auth"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

func CreateRide() gin.HandlerFunc {
	meta := map[string]string{"app": "CreateRide"}
	return func(ctx *gin.Context) {
		var ginErr *gin.Error
		vehicle := &ReqVeh{}
		r := ctx.Request
		ctx2 := db.GinToContextWithDB(ctx)
		err := bunny.Atomic(ctx2, func(ctx2 context.Context) error {

			if err := rest.UnmarshalJSONRequest(&vehicle, r); err != nil {
				err, ginErr = errors.New(ctx, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			vID := vehicle.VehicleID
			ok, err := CheckVehicle(ctx2, vID)
			if err != nil {
				err, ginErr = errors.New(ctx, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			if !ok {
				err, ginErr = errors.New(ctx, "vehicle not found", gin.ErrorTypePrivate, vehicle)
				return err
			}

			user, err := auth.FromContext(ctx)
			if err != nil {
				_, ginErr = errors.New(ctx, err.Error(), gin.ErrorTypePrivate, meta)
				errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)

			}

			//Call Ride Reby endpoint
			req, err := http.NewRequest("POST", reby.RebyHost+reby.RebyRide, nil)
			if err != nil {
				_, ginErr = errors.New(ctx, "ride call request creation failed", gin.ErrorTypePrivate, meta)
				errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
			} else {
				req.Header = reby.SetHeaders(vID.String(), reby.BearerETSEIB)
			}
			client := http.DefaultClient
			_, err = client.Do(req)
			now := time.Now()
			if err != nil {
				_, ginErr = errors.New(ctx, "ride call to the scooter failed", gin.ErrorTypePrivate, meta)
				errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
			}

			//Adding ride and updating vehicle

			ride := &models.Ride{
				ID:        models.NewRideID(),
				VehicleID: vID,
				UserID:    user.ID,
				StartedAt: now,
			}
			err = ride.Insert(ctx2)
			if err != nil {
				_, ginErr = errors.New(ctx, "ride not inserted", gin.ErrorTypePrivate, meta)
				return err
			}
			veh, err := models.FindVehicle(ctx2, vID)

			if err != nil {
				_, ginErr = errors.New(ctx, "vehicle not found in db", gin.ErrorTypePrivate, vID, meta)
				return err
			}

			veh.CurrentRideID = models.NullRideIDFrom(ride.ID)
			veh.CurrentUserID = models.NullUserIDFrom(user.ID)
			if err = veh.Update(ctx2); err != nil {
				_, ginErr = errors.New(ctx, "vehicle not updated", gin.ErrorTypePrivate, vID, meta)
				return err
			}
			writters.JsonResponse(ctx, gin.H{"message": "Ride created successfully", "ride": ride}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
	}
}

func FinishRide() gin.HandlerFunc {
	meta := map[string]string{"app": "FinishRide"}
	return func(ctx *gin.Context) {
		var ginErr *gin.Error
		ride := &ReqRide{}
		r := ctx.Request
		//Adding ride and updating vehicle
		ctx2 := db.GinToContextWithDB(ctx)
		err := bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&ride, r); err != nil {
				err, ginErr = errors.New(ctx, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			r, err := models.FindRide(ctx2, ride.ID)
			if err != nil {
				_, ginErr = errors.New(ctx, "ride not found", gin.ErrorTypePrivate, meta)
				errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
			}
			path, distance, err := CalcPath(ctx2, r)
			if err != nil {
				_, ginErr = errors.New(ctx, "path not calculated", gin.ErrorTypePrivate, meta)
				return err
			}
			r.Path.Valid = true
			r.Path.LineStringM = path
			r.Distance = distance
			r.FinishedAt = _import00.Time{
				Time:  time.Now(),
				Valid: true,
			}
			r.Duration = int32(r.FinishedAt.Time.Sub(r.StartedAt).Minutes())
			err = r.Update(ctx2)
			if err != nil {
				err, ginErr = errors.New(ctx, "ride not updated", gin.ErrorTypePrivate, meta)
				return err
			}
			veh, err := models.FindVehicle(ctx2, r.VehicleID)
			if err != nil {
				err, ginErr = errors.New(ctx, "vehicle not found", gin.ErrorTypePrivate, meta)
				return err
			}
			veh.LastRideID = veh.CurrentRideID
			veh.LastUserID = veh.CurrentUserID
			veh.CurrentRideID = models.NullRideID{}
			veh.CurrentUserID = models.NullUserID{}
			if err = veh.Update(ctx2); err != nil {
				_, ginErr = errors.New(ctx, "vehicle not updated", gin.ErrorTypePrivate, veh.ID, meta)
				return err
			}
			writters.JsonResponse(ctx, gin.H{"message": "Ride finished successfully", "ride": ride}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctx, ginErr, http.StatusBadRequest)
		}
	}
}

func AdminGetRides() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		var ginError *gin.Error
		var err error
		ctx := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx, func(ctx context.Context) error {
			rides, err := models.Rides().All(ctx)
			if err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "Rides fetched successfully", "rides": rides}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginError, http.StatusBadRequest)
		}
	}
}
