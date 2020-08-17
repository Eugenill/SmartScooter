package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/errors"
	"log"
)

var ErrBadRequest = errors.New("Bad Request")

func New(msg string) error {
	return errors.New(msg)
}

func PanicError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func ErrJsonResponse(ctx *gin.Context, err error, code int) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"error": err.Error(),
	})
}
