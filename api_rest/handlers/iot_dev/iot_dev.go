package iot_dev

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"net/http"
)

func AddIotDev() gin.HandlerFunc {
	meta := map[string]string{"app": "AddIotDev"}
	return func(ctxGin *gin.Context) {
		ginErr := &gin.Error{}
		var err error
		r := ctxGin.Request
		i := &CreateIotDev{}
		ctx2 := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&i, r); err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			iotDev := &models.IotDevice{
				ID:   models.NewIotDeviceID(),
				Name: i.Name,
			}
			err = iotDev.Insert(ctx2)
			if err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate, meta)
				return err
			}
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginErr, http.StatusBadRequest)
		}
	}
}
