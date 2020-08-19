package writters

import "github.com/gin-gonic/gin"

func JsonResponse(ctx *gin.Context, object interface{}, code int) {
	ctx.JSON(code, object)
}
