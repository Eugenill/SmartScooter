package api

import (
	"github.com/gin-gonic/gin"
)

type HTTPHandlerDefinition struct {
	Method     string
	Endpoint   string
	Handler    gin.HandlerFunc
	RouterGroup *gin.RouterGroup
}
