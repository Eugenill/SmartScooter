package handlers

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
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
)

func GetVehicles(name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writters.JsonResponse(ctx, name, 200)
	}
}

type userCreate struct {
	Username     string `json:"username" `
	Secret       string `json:"secret" `
	ContactEmail string `json:"email" `
	Admin        bool   `json:"admin" `
}

func CreateUser() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		var usr userCreate
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&usr, r); err != nil {
				return err
			}
			existUser, err := models.Users(
				qm.Where("username = ?", usr.Username),
			).Exists(ctx)
			if existUser {
				errors.ErrJsonResponse(ctxGin, errors.New("This User already exists"), 200)
				return nil
			}
			if !existUser {
				if err != nil {
					return err
				}
				o := models.User{
					ID:           models.NewUserID(),
					Username:     usr.Username,
					Secret:       usr.Secret,
					ContactEmail: usr.ContactEmail,
					Admin:        usr.Admin,
					IsDeleted:    false,
					DeletedAt:    _import00.Time{},
				}
				if err = o.Insert(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, "User created successfully", 200)
				return nil
			} else {
				return err
			}
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, err, 400)
		}
	}
}
