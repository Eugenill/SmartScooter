package user

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
	"net/http"
	"time"
)

func AdminCreateUser() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		var usr User
		var ginError *gin.Error
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&usr, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			existUser, err := models.Users(
				qm.Where("username = ?", usr.Username),
			).Exists(ctx)
			if existUser {
				err, ginError = errors.New(ctxGin, "this user already exist", gin.ErrorTypePrivate)
				return err
			}
			if !existUser {
				if err != nil {
					return err
				}
				if usr.ContactEmail == "" {
					err, ginError = errors.New(ctxGin, "please add a valid email", gin.ErrorTypePrivate)
					return err
				}
				if usr.PhoneNumber == "" {
					err, ginError = errors.New(ctxGin, "please add a valid phone number", gin.ErrorTypePrivate)
					return err
				}
				o := models.User{
					ID:           models.NewUserID(),
					Username:     usr.Username,
					Secret:       usr.Secret,
					ContactEmail: usr.ContactEmail,
					Admin:        usr.Admin,
					CreatedAt:    time.Now(),
					IsDeleted:    false,
					DeletedAt:    _import00.Time{},
				}
				if err = o.Insert(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, gin.H{"message": "User created successfully", "user": o}, http.StatusOK)
			} else if err != nil {
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

func AdminDeleteUser() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		var usr User
		var ginError *gin.Error
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&usr, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			user, err := models.Users(
				qm.Where("id = ?", usr.ID),
			).One(ctx)
			if bunny.IsErrNoRows(err) {
				err, ginError = errors.New(ctxGin, "this user does not exist", gin.ErrorTypePrivate)
				return err
			} else if err == nil {
				user.IsDeleted = true
				user.DeletedAt = _import00.Time{
					Time:  time.Now(),
					Valid: true,
				}
				if err = user.Update(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, gin.H{"message": "User deleted successfully", "user": user}, http.StatusOK)
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
func AdminEditUser() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		var usr User
		var ginError *gin.Error
		ctx := db.GinToContextWithDB(ctxGin)
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&usr, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			user, err := models.Users(
				qm.Where("id = ?", usr.ID),
			).One(ctx)
			if bunny.IsErrNoRows(err) {
				err, ginError = errors.New(ctxGin, "this user does not exist", gin.ErrorTypePrivate)
				return err
			} else if err == nil {
				if user.IsDeleted {
					err, ginError = errors.New(ctxGin, "this user has been deleted", gin.ErrorTypePrivate)
					return err
				}
				user.Username = usr.Username
				user.Secret = usr.Secret
				user.ContactEmail = usr.ContactEmail
				user.PhoneNumber = usr.PhoneNumber
				user.Admin = usr.Admin

				if err = user.Update(ctx); err != nil {
					return err
				}
				writters.JsonResponse(ctxGin, gin.H{"message": "User edited successfully", "user": user}, http.StatusOK)
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

func AdminGetUsers() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		r := ctxGin.Request
		var users Usernames
		var userResp []*RespUser
		var ginError *gin.Error
		var err error
		ctx := db.GinToContextWithDB(ctxGin)
		err = bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&users, r); err != nil {
				err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
				return err
			}
			if len(users.Usernames) != 0 {
				for _, username := range users.Usernames {
					user, err := models.Users(
						qm.Where("username = ?", username),
					).One(ctx)
					if err != nil && !bunny.IsErrNoRows(err) {
						err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
						return err
					} else if err == nil {
						userResp = append(userResp, &RespUser{
							ID:          user.ID,
							Username:    user.Username,
							Secret:      user.Secret,
							PhoneNumber: user.PhoneNumber,
							Email:       user.ContactEmail,
							Admin:       user.Admin,
							CreatedAt:   user.CreatedAt,
							IsDeleted:   user.IsDeleted,
							DeletedAt:   user.DeletedAt,
						})
					}
				}
				if len(userResp) == 0 {
					err, ginError = errors.New(ctxGin, "no users found with the given usernames", gin.ErrorTypePrivate, userResp)
					return err
				}
			} else {
				users, err := models.Users().All(ctx)
				if err != nil {
					err, ginError = errors.New(ctxGin, err.Error(), gin.ErrorTypePrivate)
					return err
				}
				for _, user := range users {
					userResp = append(userResp, &RespUser{
						ID:          user.ID,
						Username:    user.Username,
						Secret:      user.Secret,
						PhoneNumber: user.PhoneNumber,
						Email:       user.ContactEmail,
						Admin:       user.Admin,
						CreatedAt:   user.CreatedAt,
						IsDeleted:   user.IsDeleted,
						DeletedAt:   user.DeletedAt,
					})
				}
			}
			writters.JsonResponse(ctxGin, gin.H{"message": "User fetched successfully", "users": userResp}, http.StatusOK)
			return nil
		})
		if err != nil {
			errors.ErrJsonResponse(ctxGin, ginError, http.StatusBadRequest)
		}
	}
}
