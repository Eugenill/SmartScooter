package helmet

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

func CreateHelmet() gin.HandlerFunc {
	meta := map[string]string{"app": "AddHelmet"}
	return func(ctxGin *gin.Context) {
		ginErr := &gin.Error{}
		var err error
		r := ctxGin.Request
		h := &ReqHelmet{}
		ctx2 := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&h, r); err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate, meta)
				return err
			}
			helmet := &models.Helmet{
				ID:   models.NewHelmetID(),
				Name: h.Name,
			}
			err = helmet.Insert(ctx2)
			if err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate, meta)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "Helmet created successfully", "helmet": helmet}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginErr, http.StatusBadRequest)
		}
	}
}
func AdminDeleteHelmet() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		hel := &ReqHelmetID{}
		var ginError *gin.Error
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&hel, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			helmet, err := models.Helmets(
				qm.Where("id = ?", hel.ID),
			).One(ctx)
			if bunny.IsErrNoRows(err) {
				err, ginError = errors.New(ctxGin, "this helmet does not exist", gin.ErrorTypePrivate)
				return err
			} else if err == nil {
				if err = helmet.Delete(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, gin.H{"message": "helmet deleted successfully", "helmet": helmet}, http.StatusOK)
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

func AdminGetHelmets() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		var ginError *gin.Error
		var err error
		ctx := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx, func(ctx context.Context) error {
			helmets, err := models.Helmets().All(ctx)
			if err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "Helmets fetched successfully", "vehicles": helmets}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginError, http.StatusBadRequest)
		}
	}
}
