package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/errors"
	"log"
)

var ErrBadRequest = errors.New("Bad Request")

func New(ctx *gin.Context, msg string, errorType gin.ErrorType, meta ...interface{}) (error, *gin.Error) {
	error := &gin.Error{
		Err:  errors.New(msg),
		Type: errorType,
		Meta: meta,
	}
	ctx.Errors = append(ctx.Errors, error)
	return error.Err, error
}

func PanicError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func ErrJsonResponse(ctx *gin.Context, err *gin.Error, code int) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"error": err.Err.Error(),
	})
}
