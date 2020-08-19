package auth

import (
	"encoding/base64"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
	"strings"
)

func FromContext(ctx *gin.Context) (*models.User, error) {
	res, ok := ctx.Get(gin.AuthUserKey)
	if !ok {
		panic("No auth in the context")
	}
	authKey, ok := res.(string)
	if !ok {
		panic("No auth in the context")
	}
	userPass, err := base64.StdEncoding.DecodeString(authKey[6:])
	if err != nil {
		return nil, err
	}
	usr := strings.Split(string(userPass), ":")
	user, err := models.Users(
		qm.Where("user = ?", usr[0]),
		qm.Where("user = ?", usr[0]),
	).One(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
