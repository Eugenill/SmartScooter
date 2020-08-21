package vehicle

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/writters"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
	"net/http"
)

func CreateVehicle() gin.HandlerFunc {
	meta := map[string]string{"app": "AddVehicle"}
	return func(ctxGin *gin.Context) {
		vehicle := &ReqVehicle{}
		r := ctxGin.Request
		ginErr := &gin.Error{}
		var err error
		ctx2 := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&vehicle, r); err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			vehicle := &models.Vehicle{
				ID:            vehicle.ID,
				CurrentRideID: models.NullRideID{},
				LastRideID:    models.NullRideID{},
				CurrentUserID: models.NullUserID{},
				LastUserID:    models.NullUserID{},
				NumberPlate:   vehicle.NumberPlate,
				HelmetID:      vehicle.HelmetID,
				IotDeviceID:   vehicle.IotDevID,
			}
			err = vehicle.Insert(ctx2)
			if err != nil {
				_, ginErr = errors.New(ctxGin, "vehicle not inserted", gin.ErrorTypePrivate, meta)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "Vehicle created successfully", "vehicle": vehicle}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginErr, http.StatusBadRequest)
		}
	}
}

func AdminDeleteVehicle() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		veh := &ReqVehicle{}
		var ginError *gin.Error
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&veh, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			vehicle, err := models.Vehicles(
				qm.Where("id = ?", veh.ID),
			).One(ctx)
			if bunny.IsErrNoRows(err) {
				err, ginError = errors.New(ctxGin, "this vehicle does not exist", gin.ErrorTypePrivate)
				return err
			} else if err == nil {
				if err = vehicle.Delete(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, gin.H{"message": "vehicle deleted successfully", "vehicle": vehicle}, http.StatusOK)
			} else {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginError, http.StatusBadRequest)
		}
	}
}

func AdminGetVehicles() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		var ginError *gin.Error
		var err error
		ctx := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx, func(ctx context.Context) error {
			vehicles, err := models.Vehicles().All(ctx)
			if err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "Vehicles fetched successfully", "vehicles": vehicles}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginError, http.StatusBadRequest)
		}
	}
}
