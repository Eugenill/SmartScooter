package auth

import (
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
)

func PublicBasicAuth() gin.HandlerFunc {
	accounts := make(gin.Accounts)
	//accounts["public"] = "1234"
	var users models.UserSlice

	ctx := &gin.Context{}
	ctx2 := db.GinToContextWithDB(ctx)

	users, err := models.Users(
		qm.Where("is_deleted = false"),
	).All(ctx2)
	if err != nil {
		errors.New(ctx, "users not found", gin.ErrorTypePrivate)
	}
	for _, user := range users {
		accounts[user.Username] = user.Secret
	}

	if len(accounts) == 0 {
		errors.New(ctx, "users not found", gin.ErrorTypePrivate)
	}
	return gin.BasicAuth(accounts)

}

func AdminBasicAuth() gin.HandlerFunc {
	accounts := make(gin.Accounts)
	//accounts["smartscooter"] = "1234"
	var users models.UserSlice

	ctx := &gin.Context{}
	ctx2 := db.GinToContextWithDB(ctx)

	users, err := models.Users(
		qm.Where("is_deleted = false"),
		qm.Where("admin = true"),
	).All(ctx2)
	if err != nil {
		errors.New(ctx, "users not found", gin.ErrorTypePrivate)
	}
	for _, user := range users {
		accounts[user.Username] = user.Secret
	}

	if len(accounts) == 0 {
		errors.New(ctx, "users not found", gin.ErrorTypePrivate)
	}
	return gin.BasicAuth(accounts)
}
