package auth

import (
	"errors"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
)

func FromContext(ctx *gin.Context) (*models.User, error) {
	username, ok := ctx.Get(gin.AuthUserKey)
	if !ok {
		panic("No user in the context")
	}
	username, ok = username.(string)
	if !ok {
		panic("No user in the context")
	}
	ctx2 := db.GinToContextWithDB(ctx)
	user, err := models.Users(
		qm.Where("username = ?", username),
	).One(ctx2)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
