package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

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
