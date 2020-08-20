package helmet

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

func AddHelmet() gin.HandlerFunc {
	meta := map[string]string{"function": "AddHelmet"}
	return func(ctxGin *gin.Context) {
		ginErr := &gin.Error{}
		var err error
		r := ctxGin.Request
		h := &CreateHelmet{}
		ctx2 := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx2, func(ctx2 context.Context) error {
			if err := rest.UnmarshalJSONRequest(&h, r); err != nil {
				err, ginErr = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
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
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginErr, http.StatusBadRequest)
		}
	}
}
