package server

import (
	"github.com/gin-gonic/gin"
)

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey int = 0

var reqid uint64

type Router struct {
	Engine      *gin.Engine
	Middlewares []gin.HandlerFunc
	PublicRoute *gin.RouterGroup
}

func (s *Router) AddMiddlewares() {
	for _, middleware := range s.Middlewares {
		s.Engine.Use(middleware)
	}
}

func (s *Router) RunServer() error {
	err := s.Engine.Run("localhost:1234")
	return err
}