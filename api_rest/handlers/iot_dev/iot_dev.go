package iot_dev

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

func CreateIotDev() gin.HandlerFunc {
	meta := map[string]string{"app": "AddIotDev"}
	return func(ctxGin *gin.Context) {
		ginErr := &gin.Error{}
		var err error
		r := ctxGin.Request
		iot := &ReqIotDev{}
		ctx2 := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&iot, r); err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			iotDev := &models.IotDevice{
				ID:   models.NewIotDeviceID(),
				Name: iot.Name,
			}
			err = iotDev.Insert(ctx2)
			if err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate, meta)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "IotDev created successfully", "iotDev": iotDev}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginErr, http.StatusBadRequest)
		}
	}
}

func AdminDeleteIotDev() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		iot := &ReqIotDevID{}
		var ginError *gin.Error
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&iot, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			iotDev, err := models.IotDevices(
				qm.Where("id = ?", iot.ID),
			).One(ctx)
			if bunny.IsErrNoRows(err) {
				err, ginError = errors.New(ctxGin, "this iotDev does not exist", gin.ErrorTypePrivate)
				return err
			} else if err == nil {
				if err = iotDev.Delete(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, gin.H{"message": "IotDev deleted successfully", "iotDev": iotDev}, http.StatusOK)
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
func AdminGetIotDevs() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		var ginError *gin.Error
		var err error
		ctx := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx, func(ctx context.Context) error {
			iotDevs, err := models.IotDevices().All(ctx)
			if err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "Iot Devices fetched successfully", "vehicles": iotDevs}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginError, http.StatusBadRequest)
		}
	}
}
