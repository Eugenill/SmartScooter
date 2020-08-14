package handlers

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/hash"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	"log"
	"net/http"
)

func GetVehicles(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		_, err := w.Write([]byte(name))
		errors.PanicError(err)
	}
}

type userCreate struct {
	Login        string `json:"username" `
	Secret       string `json:"secret" `
	ContactEmail string `json:"email" `
}

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := bunny.ContextWithDB(r.Context(), db.DB)
		var usr userCreate
		err := bunny.Atomic(ctx, func(ctx context.Context) error {
			if err := rest.UnmarshalJSONRequest(&usr, r); err != nil {
				return err
			}
			_, err := models.Users(
				qm.Where("login = ?"),
			).One(ctx)
			if bunny.IsErrNoRows(err) {
				_, errW := w.Write([]byte("This user already exists"))
				return errW
			}
			secretHash, err := hash.HashPassword(usr.Secret)
			if err != nil {
				return err
			}
			o := models.User{
				ID:           models.NewUserID(),
				Login:        usr.Login,
				SecretHash:   secretHash,
				ContactEmail: usr.ContactEmail,
				IsDeleted:    false,
				DeletedAt:    _import00.Time{},
			}
			return o.Insert(ctx)
		})
		errors.PanicError(err)
		_, err = w.Write([]byte("User created successfully"))
		errors.PanicError(err)
	}
}
