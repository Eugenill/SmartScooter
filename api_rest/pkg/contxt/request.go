package contxt

import "github.com/gin-gonic/gin"

func RequestHeader(ctx *gin.Context, header string) string {
	return ctx.Request.Header.Get(header)
}
